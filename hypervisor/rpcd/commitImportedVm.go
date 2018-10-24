package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) CommitImportedVm(conn *srpc.Conn,
	request hypervisor.CommitImportedVmRequest,
	reply *hypervisor.CommitImportedVmResponse) error {
	*reply = hypervisor.CommitImportedVmResponse{
		errors.ErrorToString(
			t.manager.CommitImportedVm(request.IpAddress,
				conn.GetAuthInformation()))}
	return nil
}
