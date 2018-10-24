package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) DiscardVmOldUserData(conn *srpc.Conn,
	request hypervisor.DiscardVmOldUserDataRequest,
	reply *hypervisor.DiscardVmOldUserDataResponse) error {
	response := hypervisor.DiscardVmOldUserDataResponse{
		errors.ErrorToString(t.manager.DiscardVmOldUserData(request.IpAddress,
			conn.GetAuthInformation()))}
	*reply = response
	return nil
}
