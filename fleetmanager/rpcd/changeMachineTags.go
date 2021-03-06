package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/fleetmanager"
)

func (t *srpcType) ChangeMachineTags(conn *srpc.Conn,
	request fleetmanager.ChangeMachineTagsRequest,
	reply *fleetmanager.ChangeMachineTagsResponse) error {
	*reply = fleetmanager.ChangeMachineTagsResponse{
		errors.ErrorToString(t.hypervisorsManager.ChangeMachineTags(
			request.Hostname, request.Tags))}
	return nil
}
