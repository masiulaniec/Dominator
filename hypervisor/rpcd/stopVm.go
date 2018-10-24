package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) StopVm(conn *srpc.Conn,
	request hypervisor.StopVmRequest, reply *hypervisor.StopVmResponse) error {
	response := hypervisor.StopVmResponse{
		errors.ErrorToString(t.manager.StopVm(request.IpAddress,
			conn.GetAuthInformation(), request.AccessToken))}
	*reply = response
	return nil
}
