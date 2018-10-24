package client

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func makeDirectory(client *srpc.Client, dirname string) error {
	request := imageserver.MakeDirectoryRequest{dirname}
	var reply imageserver.MakeDirectoryResponse
	return client.RequestReply("ImageServer.MakeDirectory", request, &reply)
}
