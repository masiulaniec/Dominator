package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/dominator"
)

func (t *rpcType) GetDefaultImage(conn *srpc.Conn,
	request dominator.GetDefaultImageRequest,
	reply *dominator.GetDefaultImageResponse) error {
	reply.ImageName = t.herd.GetDefaultImage()
	return nil
}
