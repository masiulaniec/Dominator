package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) SnapshotVm(conn *srpc.Conn,
	request hypervisor.SnapshotVmRequest,
	reply *hypervisor.SnapshotVmResponse) error {
	err := t.manager.SnapshotVm(request.IpAddress, conn.GetAuthInformation(),
		request.ForceIfNotStopped, request.RootOnly)
	*reply = hypervisor.SnapshotVmResponse{errors.ErrorToString(err)}
	return nil
}
