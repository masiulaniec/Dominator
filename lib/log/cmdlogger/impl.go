package cmdlogger

import (
	"flag"
	"log"

	"github.com/masiulaniec/Dominator/lib/log/debuglogger"
)

func init() {
	flag.BoolVar(&stdOptions.Datestamps, "logDatestamps", false,
		"If true, prefix logs with datestamps")
	flag.IntVar(&stdOptions.DebugLevel, "logDebugLevel", -1, "Debug log level")
}

func newLogger(options Options) *debuglogger.Logger {
	if options.DebugLevel < -1 {
		options.DebugLevel = -1
	}
	if options.DebugLevel > 65535 {
		options.DebugLevel = 65535
	}
	logFlags := 0
	if options.Datestamps {
		logFlags |= log.LstdFlags
	}
	logger := debuglogger.New(log.New(options.Writer, "", logFlags))
	logger.SetLevel(int16(options.DebugLevel))
	return logger
}
