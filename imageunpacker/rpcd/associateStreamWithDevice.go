package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) AssociateStreamWithDevice(conn *srpc.Conn,
	request proto.AssociateStreamWithDeviceRequest,
	reply *proto.AssociateStreamWithDeviceResponse) error {
	return t.unpacker.AssociateStreamWithDevice(request.StreamName,
		request.DeviceId)
}
