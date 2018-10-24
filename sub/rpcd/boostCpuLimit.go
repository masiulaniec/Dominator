package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
)

func (t *rpcType) BoostCpuLimit(conn *srpc.Conn,
	request sub.BoostCpuLimitRequest, reply *sub.BoostCpuLimitResponse) error {
	t.scannerConfiguration.BoostCpuLimit(t.logger)
	return nil
}
