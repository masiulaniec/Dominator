package client

import (
	"github.com/masiulaniec/Dominator/lib/image"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func addImage(client *srpc.Client, name string, img *image.Image) error {
	request := imageserver.AddImageRequest{name, img}
	var reply imageserver.AddImageResponse
	return client.RequestReply("ImageServer.AddImage", request, &reply)
}

func addImageTrusted(client *srpc.Client, name string, img *image.Image) error {
	request := imageserver.AddImageRequest{name, img}
	var reply imageserver.AddImageResponse
	return client.RequestReply("ImageServer.AddImageTrusted", request, &reply)
}
