package memory

import (
	"github.com/masiulaniec/Dominator/lib/hash"
)

func newObjectServer() *ObjectServer {
	objSrv := new(ObjectServer)
	objSrv.objectMap = make(map[hash.Hash][]byte)
	return objSrv
}
