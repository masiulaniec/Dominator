package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	proto "github.com/masiulaniec/Dominator/proto/fleetmanager"
)

func (t *srpcType) ListHypervisorsInLocation(conn *srpc.Conn,
	request proto.ListHypervisorsInLocationRequest,
	reply *proto.ListHypervisorsInLocationResponse) error {
	addresses, err := t.hypervisorsManager.ListHypervisorsInLocation(request)
	*reply = proto.ListHypervisorsInLocationResponse{
		HypervisorAddresses: addresses,
		Error:               errors.ErrorToString(err),
	}
	return nil
}
