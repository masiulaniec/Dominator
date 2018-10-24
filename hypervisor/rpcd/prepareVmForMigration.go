package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) PrepareVmForMigration(conn *srpc.Conn,
	request hypervisor.PrepareVmForMigrationRequest,
	reply *hypervisor.PrepareVmForMigrationResponse) error {
	*reply = hypervisor.PrepareVmForMigrationResponse{
		Error: errors.ErrorToString(
			t.manager.PrepareVmForMigration(request.IpAddress,
				conn.GetAuthInformation(), request.AccessToken,
				request.Enable)),
	}
	return nil
}
