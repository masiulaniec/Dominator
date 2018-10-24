package rpcd

import (
	"github.com/masiulaniec/Dominator/dom/herd"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/srpc"
)

type rpcType struct {
	herd   *herd.Herd
	logger log.Logger
}

func Setup(herd *herd.Herd, logger log.Logger) {
	rpcObj := &rpcType{
		herd:   herd,
		logger: logger}
	srpc.RegisterName("Dominator", rpcObj)
}
