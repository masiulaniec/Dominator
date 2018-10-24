package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/dominator"
	"github.com/masiulaniec/Dominator/proto/sub"
)

func (t *rpcType) ConfigureSubs(conn *srpc.Conn,
	request dominator.ConfigureSubsRequest,
	reply *dominator.ConfigureSubsResponse) error {
	if conn.Username() == "" {
		t.logger.Printf("ConfigureSubs()\n")
	} else {
		t.logger.Printf("ConfigureSubs(): by %s\n", conn.Username())
	}
	return t.herd.ConfigureSubs(sub.Configuration(request))
}
