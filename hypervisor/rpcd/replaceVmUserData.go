package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) ReplaceVmUserData(conn *srpc.Conn,
	request hypervisor.ReplaceVmUserDataRequest,
	reply *hypervisor.ReplaceVmUserDataResponse) error {
	response := hypervisor.ReplaceVmUserDataResponse{
		errors.ErrorToString(t.manager.ReplaceVmUserData(request.IpAddress,
			conn, request.Size, conn.GetAuthInformation()))}
	*reply = response
	return nil
}
