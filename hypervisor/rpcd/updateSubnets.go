package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) UpdateSubnets(conn *srpc.Conn,
	request hypervisor.UpdateSubnetsRequest,
	reply *hypervisor.UpdateSubnetsResponse) error {
	*reply = hypervisor.UpdateSubnetsResponse{
		errors.ErrorToString(t.manager.UpdateSubnets(request))}
	return nil
}
