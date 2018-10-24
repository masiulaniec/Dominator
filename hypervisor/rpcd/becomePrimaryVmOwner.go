package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) BecomePrimaryVmOwner(conn *srpc.Conn,
	request hypervisor.BecomePrimaryVmOwnerRequest,
	reply *hypervisor.BecomePrimaryVmOwnerResponse) error {
	*reply = hypervisor.BecomePrimaryVmOwnerResponse{
		errors.ErrorToString(
			t.manager.BecomePrimaryVmOwner(request.IpAddress,
				conn.GetAuthInformation()))}
	return nil
}
