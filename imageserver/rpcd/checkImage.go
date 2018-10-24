package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func (t *srpcType) CheckImage(conn *srpc.Conn,
	request imageserver.CheckImageRequest,
	reply *imageserver.CheckImageResponse) error {
	var response imageserver.CheckImageResponse
	response.ImageExists = t.imageDataBase.CheckImage(request.ImageName)
	*reply = response
	return nil
}
