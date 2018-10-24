package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/objectserver"
)

func (t *srpcType) CheckObjects(conn *srpc.Conn,
	request objectserver.CheckObjectsRequest,
	reply *objectserver.CheckObjectsResponse) error {
	sizes, err := t.objectServer.CheckObjects(request.Hashes)
	if err != nil {
		_, err = conn.WriteString(err.Error() + "\n")
		return err
	}
	reply.ObjectSizes = sizes
	return nil
}
