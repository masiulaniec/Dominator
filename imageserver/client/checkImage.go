package client

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func checkImage(client *srpc.Client, name string) (bool, error) {
	request := imageserver.CheckImageRequest{name}
	var reply imageserver.CheckImageResponse
	err := client.RequestReply("ImageServer.CheckImage", request, &reply)
	if err != nil {
		return false, err
	}
	return reply.ImageExists, nil
}
