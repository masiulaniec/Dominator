package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/srpc"
)

func (t *srpcType) GetVmVolume(conn *srpc.Conn, decoder srpc.Decoder,
	encoder srpc.Encoder) error {
	return t.manager.GetVmVolume(conn, decoder, encoder)
}
