package client

import (
	"time"

	"github.com/masiulaniec/Dominator/lib/image"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func getImage(client *srpc.Client, name string, timeout time.Duration) (
	*image.Image, error) {
	request := imageserver.GetImageRequest{ImageName: name, Timeout: timeout}
	var reply imageserver.GetImageResponse
	err := client.RequestReply("ImageServer.GetImage", request, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Image, nil
}
