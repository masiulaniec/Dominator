package rpcd

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/masiulaniec/Dominator/lib/hash"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/objectserver"
)

var exclusive sync.RWMutex

func (objSrv *srpcType) GetObjects(conn *srpc.Conn, decoder srpc.Decoder,
	encoder srpc.Encoder) error {
	defer conn.Flush()
	var request objectserver.GetObjectsRequest
	var response objectserver.GetObjectsResponse
	if request.Exclusive {
		exclusive.Lock()
		defer exclusive.Unlock()
	} else {
		exclusive.RLock()
		defer exclusive.RUnlock()
		objSrv.getSemaphore <- true
		defer releaseSemaphore(objSrv.getSemaphore)
	}
	var err error
	if err = decoder.Decode(&request); err != nil {
		response.ResponseString = err.Error()
		return encoder.Encode(response)
	}
	response.ObjectSizes, err = objSrv.objectServer.CheckObjects(request.Hashes)
	if err != nil {
		response.ResponseString = err.Error()
		return encoder.Encode(response)
	}
	// First a quick check for existence. If any objects missing, fail request.
	var firstMissingObject *hash.Hash
	numMissingObjects := 0
	for index, hashVal := range request.Hashes {
		if response.ObjectSizes[index] < 1 {
			firstMissingObject = &hashVal
			numMissingObjects++
		}
	}
	if firstMissingObject != nil {
		if numMissingObjects == 1 {
			response.ResponseString = fmt.Sprintf("unknown object: %x",
				*firstMissingObject)
		} else {
			response.ResponseString = fmt.Sprintf(
				"first of %d unknown objects: %x", numMissingObjects,
				*firstMissingObject)
		}
		return encoder.Encode(response)
	}
	objectsReader, err := objSrv.objectServer.GetObjects(request.Hashes)
	if err != nil {
		response.ResponseString = err.Error()
		return encoder.Encode(response)
	}
	defer objectsReader.Close()
	if err := encoder.Encode(response); err != nil {
		return err
	}
	conn.Flush()
	for _, hashVal := range request.Hashes {
		length, reader, err := objectsReader.NextObject()
		if err != nil {
			objSrv.logger.Println(err)
			return err
		}
		nCopied, err := io.Copy(conn.Writer, reader)
		reader.Close()
		if err != nil {
			objSrv.logger.Printf("Error copying: %s\n", err)
			return err
		}
		if nCopied != int64(length) {
			txt := fmt.Sprintf("Expected length: %d, got: %d for: %x",
				length, nCopied, hashVal)
			objSrv.logger.Printf(txt)
			return errors.New(txt)
		}
	}
	objSrv.logger.Debugf(0, "GetObjects() sent: %d objects\n",
		len(request.Hashes))
	return nil
}

func releaseSemaphore(semaphore <-chan bool) {
	<-semaphore
}
