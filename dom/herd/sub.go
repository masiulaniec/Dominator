package herd

import (
	"fmt"
	"github.com/Symantec/Dominator/lib/constants"
	"github.com/Symantec/Dominator/lib/hash"
	subproto "github.com/Symantec/Dominator/proto/sub"
	"net/rpc"
	"strings"
	"time"
)

func (sub *Sub) tryMakeBusy() bool {
	sub.busyMutex.Lock()
	defer sub.busyMutex.Unlock()
	if sub.busy {
		return false
	}
	sub.busy = true
	return true
}

func (sub *Sub) makeUnbusy() {
	sub.busyMutex.Lock()
	defer sub.busyMutex.Unlock()
	sub.busy = false
}

func (sub *Sub) connectAndPoll() {
	sub.status = statusConnecting
	hostname := strings.SplitN(sub.hostname, "*", 2)[0]
	var err error
	connection, err := rpc.DialHTTP("tcp",
		fmt.Sprintf("%s:%d", hostname, constants.SubPortNumber))
	if err != nil {
		sub.status = statusFailedToConnect
		return
	}
	defer connection.Close()
	sub.status = statusWaitingToPoll
	sub.herd.pollSemaphore <- true
	sub.status = statusPolling
	sub.poll(connection)
	<-sub.herd.pollSemaphore
}

func (sub *Sub) poll(connection *rpc.Client) {
	var request subproto.PollRequest
	request.HaveGeneration = sub.generationCount
	var reply subproto.PollResponse
	err := connection.Call("Subd.Poll", request, &reply)
	if err != nil {
		sub.status = statusFailedToPoll
		fmt.Printf("Error calling\t%s\n", err)
		return
	}
	sub.lastSuccessfulPoll = time.Now()
	if reply.GenerationCount == 0 {
		sub.fileSystem = nil
	}
	if fs := reply.FileSystem; fs != nil {
		fs.RebuildPointers()
		sub.fileSystem = fs
		sub.generationCount = reply.GenerationCount
		// TODO(rgooch): Remove debugging output.
		fmt.Printf("Polled: %s, GenerationCount=%d\n",
			sub.hostname, reply.GenerationCount)
	}
	if sub.fileSystem == nil {
		sub.status = statusSubNotReady
		return
	}
	if reply.FetchInProgress {
		sub.status = statusFetching
		return
	}
	if reply.UpdateInProgress {
		sub.status = statusUpdating
		return
	}
	if sub.generationCountAtChangeStart == sub.generationCount {
		sub.status = statusWaitingForNextPoll
		return
	}
	if idle, status := sub.fetchMissingObjects(connection,
		sub.requiredImage); !idle {
		sub.status = status
		return
	}
	if idle, status := sub.sendUpdate(connection); !idle {
		sub.status = status
		return
	}
	if idle, status := sub.fetchMissingObjects(connection,
		sub.plannedImage); !idle {
		sub.status = status
		return
	}
	sub.status = statusSynced
}

// Returns true if all required objects are available.
func (sub *Sub) fetchMissingObjects(connection *rpc.Client, imageName string) (
	bool, uint) {
	if imageName == "" {
		return false, statusImageNotReady
	}
	image := sub.herd.getImage(imageName)
	if image == nil {
		return false, statusImageNotReady
	}
	missingObjects := make(map[hash.Hash]bool)
	for _, inode := range image.FileSystem.RegularInodeTable {
		if inode.Size > 0 {
			missingObjects[inode.Hash] = true
		}
	}
	for _, hash := range sub.fileSystem.ObjectCache {
		delete(missingObjects, hash)
	}
	for _, inode := range sub.fileSystem.RegularInodeTable {
		if inode.Size > 0 {
			delete(missingObjects, inode.Hash)
		}
	}
	if len(missingObjects) < 1 {
		return true, statusSynced
	}
	// TODO(rgooch): Remove debugging output.
	fmt.Printf("Objects needing to be fetched: %d\n", len(missingObjects))
	var request subproto.FetchRequest
	var reply subproto.FetchResponse
	request.ServerAddress = sub.herd.imageServerAddress
	for hash, _ := range missingObjects {
		request.Hashes = append(request.Hashes, hash)
	}
	err := connection.Call("Subd.Fetch", request, &reply)
	if err != nil {
		fmt.Printf("Error calling\t%s\n", err)
		return false, statusFailedToFetch
	}
	sub.generationCountAtChangeStart = sub.generationCount
	return false, statusFetching
}

// Returns true if no update needs to be performed.
func (sub *Sub) sendUpdate(connection *rpc.Client) (bool, uint) {
	// TODO(rgooch): Implement this.
	return false, statusUpdating
}