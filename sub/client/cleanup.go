package client

import (
	"github.com/masiulaniec/Dominator/lib/hash"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
)

func cleanup(client *srpc.Client, hashes []hash.Hash) error {
	request := sub.CleanupRequest{hashes}
	var reply sub.CleanupResponse
	return client.RequestReply("Subd.Cleanup", request, &reply)
}
