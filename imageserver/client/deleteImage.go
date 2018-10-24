package client

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func deleteImage(client *srpc.Client, name string) error {
	request := imageserver.DeleteImageRequest{name}
	var reply imageserver.DeleteImageResponse
	return client.RequestReply("ImageServer.DeleteImage", request, &reply)
}
