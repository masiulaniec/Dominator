package client

import (
	"github.com/masiulaniec/Dominator/lib/hash"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
)

func fetch(client *srpc.Client, serverAddress string,
	hashes []hash.Hash) error {
	request := sub.FetchRequest{ServerAddress: serverAddress, Hashes: hashes}
	var reply sub.FetchResponse
	return client.RequestReply("Subd.Fetch", request, &reply)
}
