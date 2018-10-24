package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) RestoreVmFromSnapshot(conn *srpc.Conn,
	request hypervisor.RestoreVmFromSnapshotRequest,
	reply *hypervisor.RestoreVmFromSnapshotResponse) error {
	response := hypervisor.RestoreVmFromSnapshotResponse{
		errors.ErrorToString(t.manager.RestoreVmFromSnapshot(request.IpAddress,
			conn.GetAuthInformation(), request.ForceIfNotStopped))}
	*reply = response
	return nil
}
