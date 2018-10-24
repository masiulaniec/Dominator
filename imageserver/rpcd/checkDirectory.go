package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func (t *srpcType) CheckDirectory(conn *srpc.Conn,
	request imageserver.CheckDirectoryRequest,
	reply *imageserver.CheckDirectoryResponse) error {
	response := imageserver.CheckDirectoryResponse{
		t.imageDataBase.CheckDirectory(request.DirectoryName)}
	*reply = response
	return nil
}
