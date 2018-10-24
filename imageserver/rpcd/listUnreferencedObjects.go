package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/imageserver"
)

func (t *srpcType) ListUnreferencedObjects(conn *srpc.Conn,
	decoder srpc.Decoder, encoder srpc.Encoder) error {
	for hashVal, size := range t.imageDataBase.ListUnreferencedObjects() {
		obj := imageserver.Object{Hash: hashVal, Size: size}
		if err := encoder.Encode(obj); err != nil {
			return err
		}
	}
	return encoder.Encode(imageserver.Object{})
}
