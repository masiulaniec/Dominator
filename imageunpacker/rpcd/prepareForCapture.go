package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) PrepareForCapture(conn *srpc.Conn,
	request proto.PrepareForCaptureRequest,
	reply *proto.PrepareForCaptureResponse) error {
	return t.unpacker.PrepareForCapture(request.StreamName)
}
