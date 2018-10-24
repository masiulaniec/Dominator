package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imageunpacker"
)

func (t *srpcType) PrepareForUnpack(conn *srpc.Conn,
	request proto.PrepareForUnpackRequest,
	reply *proto.PrepareForUnpackResponse) error {
	return t.unpacker.PrepareForUnpack(request.StreamName,
		request.SkipIfPrepared, request.DoNotWaitForResult)
}
