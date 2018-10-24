package filegen

import (
	"bytes"
	"time"

	"github.com/masiulaniec/Dominator/lib/json"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/mdb"
)

type jsonType struct{}

func (m *Manager) registerMdbGeneratorForPath(pathname string) {
	close(m.registerDataGeneratorForPath(pathname, jsonType{}))
}

func (jsonType) Generate(machine mdb.Machine, logger log.Logger) (
	[]byte, time.Time, error) {
	buffer := new(bytes.Buffer)
	if err := json.WriteWithIndent(buffer, "    ", machine); err != nil {
		return nil, time.Time{}, err
	}
	buffer.WriteString("\n")
	return buffer.Bytes(), time.Time{}, nil
}
