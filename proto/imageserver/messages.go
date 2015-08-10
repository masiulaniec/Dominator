package imageserver

import (
	"github.com/Symantec/Dominator/lib/hash"
	"github.com/Symantec/Dominator/lib/image"
	"github.com/Symantec/Dominator/proto/common"
)

type addFile struct {
	ObjectData   []byte
	ExpectedHash *hash.Hash
}

type AddFilesRequest struct {
	ObjectsToAdd []*addFile
}

type AddFilesResponse struct {
	Hashes []hash.Hash
}

type AddImageRequest struct {
	ImageName string
	Image     *image.Image
}

type AddImageResponse struct {
}

type CheckImageRequest struct {
	ImageName string
}

type CheckImageResponse struct {
	ImageExists bool
}

type DeleteImageRequest struct {
	ImageName string
}

type DeleteImageResponse common.StatusResponse

type GetFilesRequest struct {
	Objects []hash.Hash
}

type GetFilesResponse struct {
	ObjectSizes []uint64
}

type GetImageRequest struct {
	ImageName string
}

type GetImageResponse struct {
	Image *image.Image
}

type ListImagesRequest struct {
}

type ListImagesResponse struct {
	ImageNames []string
}
