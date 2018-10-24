package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) UnpackImage(conn *srpc.Conn,
	request proto.UnpackImageRequest,
	reply *proto.UnpackImageResponse) error {
	return t.unpacker.UnpackImage(request.StreamName, request.ImageLeafName)
}
