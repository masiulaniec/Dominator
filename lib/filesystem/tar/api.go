package tar

import (
	"archive/tar"
	"io"

	"github.com/masiulaniec/Dominator/lib/filesystem"
	"github.com/masiulaniec/Dominator/lib/objectserver"
)

func Encode(tarWriter *tar.Writer, fileSystem *filesystem.FileSystem,
	objectsGetter objectserver.ObjectsGetter) error {
	return encode(tarWriter, fileSystem, objectsGetter)
}

func Write(writer io.Writer, fileSystem *filesystem.FileSystem,
	objectsGetter objectserver.ObjectsGetter) error {
	return write(writer, fileSystem, objectsGetter)
}
