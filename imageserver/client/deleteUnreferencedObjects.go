package client

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func deleteUnreferencedObjects(client *srpc.Client, percentage uint8,
	bytes uint64) error {
	request := imageserver.DeleteUnreferencedObjectsRequest{percentage, bytes}
	var reply imageserver.DeleteUnreferencedObjectsResponse
	return client.RequestReply("ImageServer.DeleteUnreferencedObjects",
		request, &reply)
}
