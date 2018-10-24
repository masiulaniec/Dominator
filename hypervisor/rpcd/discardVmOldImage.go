package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) DiscardVmOldImage(conn *srpc.Conn,
	request hypervisor.DiscardVmOldImageRequest,
	reply *hypervisor.DiscardVmOldImageResponse) error {
	response := hypervisor.DiscardVmOldImageResponse{
		errors.ErrorToString(t.manager.DiscardVmOldImage(request.IpAddress,
			conn.GetAuthInformation()))}
	*reply = response
	return nil
}
