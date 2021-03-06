package client

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
)

func setConfiguration(client *srpc.Client, config sub.Configuration) error {
	var request sub.SetConfigurationRequest
	request = sub.SetConfigurationRequest(config)
	var reply sub.SetConfigurationResponse
	return client.RequestReply("Subd.SetConfiguration", request, &reply)
}
