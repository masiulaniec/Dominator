package rpcd

import (
	"bytes"
	"io"
	"os"
	"path"
	"syscall"

	"github.com/masiulaniec/Dominator/lib/fsutil"
	"github.com/masiulaniec/Dominator/lib/hash"
	"github.com/masiulaniec/Dominator/lib/objectcache"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/objectserver/rpcd/lib"
)

const (
	dirPerms = syscall.S_IRWXU
)

type objectServer struct {
	baseDir string
}

func (t *addObjectsHandlerType) AddObjects(conn *srpc.Conn,
	decoder srpc.Decoder, encoder srpc.Encoder) error {
	defer t.scannerConfiguration.BoostCpuLimit(t.logger)
	objSrv := &objectServer{t.objectsDir}
	return lib.AddObjects(conn, decoder, encoder, objSrv, t.logger)
}

func (objSrv *objectServer) AddObject(reader io.Reader, length uint64,
	expectedHash *hash.Hash) (hash.Hash, bool, error) {
	hashVal, data, err := objectcache.ReadObject(reader, length, expectedHash)
	if err != nil {
		return hashVal, false, err
	}
	filename := path.Join(objSrv.baseDir, objectcache.HashToFilename(hashVal))
	if err = os.MkdirAll(path.Dir(filename), dirPerms); err != nil {
		return hashVal, false, err
	}
	if err := fsutil.CopyToFile(filename, filePerms, bytes.NewReader(data),
		length); err != nil {
		return hashVal, false, err
	}
	return hashVal, true, nil
}
