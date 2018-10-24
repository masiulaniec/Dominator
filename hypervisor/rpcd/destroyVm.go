package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) DestroyVm(conn *srpc.Conn,
	request hypervisor.DestroyVmRequest,
	reply *hypervisor.DestroyVmResponse) error {
	response := hypervisor.DestroyVmResponse{
		errors.ErrorToString(t.manager.DestroyVm(request.IpAddress,
			conn.GetAuthInformation(), request.AccessToken))}
	*reply = response
	return nil
}
