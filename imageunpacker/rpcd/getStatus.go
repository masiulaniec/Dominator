package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) GetStatus(conn *srpc.Conn, request proto.GetStatusRequest,
	reply *proto.GetStatusResponse) error {
	*reply = t.unpacker.GetStatus()
	return nil
}
