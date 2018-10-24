package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) ChangeVmOwnerUsers(conn *srpc.Conn,
	request hypervisor.ChangeVmOwnerUsersRequest,
	reply *hypervisor.ChangeVmOwnerUsersResponse) error {
	response := hypervisor.ChangeVmOwnerUsersResponse{
		errors.ErrorToString(
			t.manager.ChangeVmOwnerUsers(request.IpAddress,
				conn.GetAuthInformation(),
				request.OwnerUsers))}
	*reply = response
	return nil
}
