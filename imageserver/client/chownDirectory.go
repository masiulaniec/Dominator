package client

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func chownDirectory(client *srpc.Client, dirname, ownerGroup string) error {
	request := imageserver.ChangeOwnerRequest{DirectoryName: dirname,
		OwnerGroup: ownerGroup}
	var reply imageserver.ChangeOwnerResponse
	return client.RequestReply("ImageServer.ChownDirectory", request, &reply)
}
