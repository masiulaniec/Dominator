package imageserver

import (
	"time"

	"github.com/masiulaniec/Dominator/lib/hash"
	"github.com/masiulaniec/Dominator/lib/image"
)

type AddImageRequest struct {
	ImageName string
	Image     *image.Image
}

type AddImageResponse struct{}

type ChangeOwnerRequest struct {
	DirectoryName string
	OwnerGroup    string
}

type ChangeOwnerResponse struct{}

type CheckDirectoryRequest struct {
	DirectoryName string
}

type CheckDirectoryResponse struct {
	DirectoryExists bool
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

type DeleteImageResponse struct{}

type DeleteUnreferencedObjectsRequest struct {
	Percentage uint8
	Bytes      uint64
}

type DeleteUnreferencedObjectsResponse struct{}

type FindLatestImageRequest struct {
	DirectoryName        string
	IgnoreExpiringImages bool
}

type FindLatestImageResponse struct {
	ImageName string
	Error     string
}

type GetImageRequest struct {
	ImageName        string
	IgnoreFilesystem bool
	Timeout          time.Duration
}

type GetImageResponse struct {
	Image *image.Image
}

const (
	OperationAddImage = iota
	OperationDeleteImage
	OperationMakeDirectory
)

// The GetImageUpdates() RPC is fully streamed.
// The client sends no information to the server.
// The server sends a stream of ImageUpdate messages.

type ImageUpdate struct {
	Name      string // "" signifies initial list is sent, changes to follow.
	Directory *image.Directory
	Operation uint
}

// The ListDirectories() RPC is fully streamed.
// The client sends no information to the server.
// The server sends a stream of image.Directory values with an empty string
// for the Name field signifying the end of the list.

// The ListImages() RPC is fully streamed.
// The client sends no information to the server.
// The server sends a stream of strings (image names) with an empty string
// signifying the end of the list.

// The ListUnreferencedObjects() RPC is fully streamed.
// The client sends no information to the server.
// The server sends a stream of Object values with a zero Size field signifying
// the end of the stream.

type Object struct {
	Hash hash.Hash
	Size uint64
}

type MakeDirectoryRequest struct {
	DirectoryName string
}

type MakeDirectoryResponse struct{}
