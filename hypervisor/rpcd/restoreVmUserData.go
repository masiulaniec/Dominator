package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) RestoreVmUserData(conn *srpc.Conn,
	request hypervisor.RestoreVmUserDataRequest,
	reply *hypervisor.RestoreVmUserDataResponse) error {
	response := hypervisor.RestoreVmUserDataResponse{
		errors.ErrorToString(t.manager.RestoreVmUserData(request.IpAddress,
			conn.GetAuthInformation()))}
	*reply = response
	return nil
}
