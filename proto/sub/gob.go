package sub

import (
	"encoding/gob"

	"github.com/masiulaniec/Dominator/lib/filesystem"
)

func init() {
	gob.Register(&filesystem.RegularInode{})
	gob.Register(&filesystem.SymlinkInode{})
	gob.Register(&filesystem.SpecialInode{})
	gob.Register(&filesystem.DirectoryInode{})
}
