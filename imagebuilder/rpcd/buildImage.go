package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/imaginator"
)

func (t *srpcType) BuildImage(conn *srpc.Conn, request proto.BuildImageRequest,
	reply *proto.BuildImageResponse) error {
	name, buildLog, err := t.builder.BuildImage(request.StreamName,
		request.ExpiresIn, request.GitBranch, request.MaxSourceAge)
	reply.ImageName = name
	reply.BuildLog = buildLog
	if err != nil {
		reply.ErrorString = err.Error()
	}
	return nil
}
