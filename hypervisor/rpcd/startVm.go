package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) StartVm(conn *srpc.Conn,
	request hypervisor.StartVmRequest,
	reply *hypervisor.StartVmResponse) error {
	dhcpTimedOut, err := t.manager.StartVm(request.IpAddress,
		conn.GetAuthInformation(), request.AccessToken, request.DhcpTimeout)
	response := hypervisor.StartVmResponse{dhcpTimedOut,
		errors.ErrorToString(err)}
	*reply = response
	return nil
}
