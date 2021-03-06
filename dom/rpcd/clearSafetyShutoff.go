package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/dominator"
)

func (t *rpcType) ClearSafetyShutoff(conn *srpc.Conn,
	request dominator.ClearSafetyShutoffRequest,
	reply *dominator.ClearSafetyShutoffResponse) error {
	if conn.Username() == "" {
		t.logger.Printf("ClearSafetyShutoff(%s)\n", request.Hostname)
	} else {
		t.logger.Printf("ClearSafetyShutoff(%s): by %s\n",
			request.Hostname, conn.Username())
	}
	return t.herd.ClearSafetyShutoff(request.Hostname)
}
