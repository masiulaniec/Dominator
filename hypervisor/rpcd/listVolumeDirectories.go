package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/hypervisor"
)

func (t *srpcType) ListVolumeDirectories(conn *srpc.Conn,
	request hypervisor.ListVolumeDirectoriesRequest,
	reply *hypervisor.ListVolumeDirectoriesResponse) error {
	directories, err := t.listVolumeDirectories(conn)
	*reply = hypervisor.ListVolumeDirectoriesResponse{directories,
		errors.ErrorToString(err)}
	return nil
}

func (t *srpcType) listVolumeDirectories(conn *srpc.Conn) ([]string, error) {
	if err := testIfLoopback(conn); err != nil {
		return nil, err
	}
	return t.manager.ListVolumeDirectories(), nil
}
