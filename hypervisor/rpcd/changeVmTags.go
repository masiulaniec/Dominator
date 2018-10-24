package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) ChangeVmTags(conn *srpc.Conn,
	request hypervisor.ChangeVmTagsRequest,
	reply *hypervisor.ChangeVmTagsResponse) error {
	response := hypervisor.ChangeVmTagsResponse{
		errors.ErrorToString(
			t.manager.ChangeVmTags(request.IpAddress, conn.GetAuthInformation(),
				request.Tags))}
	*reply = response
	return nil
}
