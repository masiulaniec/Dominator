package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imagepublishers/amipublisher"
	"github.com/masiulaniec/Dominator/lib/log"
)

func stopIdleUnpackersSubcommand(args []string, logger log.DebugLogger) {
	err := amipublisher.StopIdleUnpackers(targets, skipTargets, *instanceName,
		*maxIdleTime, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stopping idle unpackers: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
