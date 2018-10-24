package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/dominator"
)

func (t *rpcType) SetDefaultImage(conn *srpc.Conn,
	request dominator.SetDefaultImageRequest,
	reply *dominator.SetDefaultImageResponse) error {
	if conn.Username() == "" {
		t.logger.Printf("SetDefaultImage(%s)\n", request.ImageName)
	} else {
		t.logger.Printf("SetDefaultImage(%s): by %s\n",
			request.ImageName, conn.Username())
	}
	return t.herd.SetDefaultImage(request.ImageName)
}
