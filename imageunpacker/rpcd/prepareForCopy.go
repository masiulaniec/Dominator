package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) PrepareForCopy(conn *srpc.Conn,
	request proto.PrepareForCopyRequest,
	reply *proto.PrepareForCopyResponse) error {
	return t.unpacker.PrepareForCopy(request.StreamName)
}
