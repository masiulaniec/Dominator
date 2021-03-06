package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/image"
	"github.com/masiulaniec/Dominator/lib/srpc"
)

func (t *srpcType) ListDirectories(conn *srpc.Conn, decoder srpc.Decoder,
	encoder srpc.Encoder) error {
	for _, directory := range t.imageDataBase.ListDirectories() {
		if err := encoder.Encode(directory); err != nil {
			return err
		}
	}
	return encoder.Encode(image.Directory{})
}
