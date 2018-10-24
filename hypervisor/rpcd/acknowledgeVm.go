package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) AcknowledgeVm(conn *srpc.Conn,
	request hypervisor.AcknowledgeVmRequest,
	reply *hypervisor.AcknowledgeVmResponse) error {
	response := hypervisor.AcknowledgeVmResponse{
		errors.ErrorToString(t.manager.AcknowledgeVm(request.IpAddress,
			conn.GetAuthInformation()))}
	*reply = response
	return nil
}
