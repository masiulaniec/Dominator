package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/dominator"
)

func (t *rpcType) GetSubsConfiguration(conn *srpc.Conn,
	request dominator.GetSubsConfigurationRequest,
	reply *dominator.GetSubsConfigurationResponse) error {
	*reply = dominator.GetSubsConfigurationResponse(
		t.herd.GetSubsConfiguration())
	return nil
}
