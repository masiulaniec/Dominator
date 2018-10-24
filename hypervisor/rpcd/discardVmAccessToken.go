package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) DiscardVmAccessToken(conn *srpc.Conn,
	request hypervisor.DiscardVmAccessTokenRequest,
	reply *hypervisor.DiscardVmAccessTokenResponse) error {
	*reply = hypervisor.DiscardVmAccessTokenResponse{
		errors.ErrorToString(t.manager.DiscardVmAccessToken(request.IpAddress,
			conn.GetAuthInformation(), request.AccessToken))}
	return nil
}
