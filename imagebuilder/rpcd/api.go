package rpcd

import (
	"io"

	"github.com/masiulaniec/Dominator/imagebuilder/builder"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/srpc"
)

type srpcType struct {
	builder *builder.Builder
	logger  log.Logger
}

type htmlWriter srpcType

func (hw *htmlWriter) WriteHtml(writer io.Writer) {
	hw.writeHtml(writer)
}

func Setup(builder *builder.Builder, logger log.Logger) (*htmlWriter, error) {
	srpcObj := &srpcType{
		builder: builder,
		logger:  logger,
	}
	srpc.RegisterName("Imaginator", srpcObj)
	return (*htmlWriter)(srpcObj), nil
}
