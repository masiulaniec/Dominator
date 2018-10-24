package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) ExportImage(conn *srpc.Conn,
	request proto.ExportImageRequest,
	reply *proto.ExportImageResponse) error {
	return t.unpacker.ExportImage(request.StreamName, request.Type,
		request.Destination)
}
