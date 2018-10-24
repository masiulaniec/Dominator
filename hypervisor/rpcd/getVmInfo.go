package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) GetVmInfo(conn *srpc.Conn,
	request hypervisor.GetVmInfoRequest,
	reply *hypervisor.GetVmInfoResponse) error {
	info, err := t.manager.GetVmInfo(request.IpAddress)
	response := hypervisor.GetVmInfoResponse{
		VmInfo: info,
		Error:  errors.ErrorToString(err),
	}
	*reply = response
	return nil
}
