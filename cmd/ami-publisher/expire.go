package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imagepublishers/amipublisher"
	"github.com/masiulaniec/Dominator/lib/log"
)

func expireSubcommand(args []string, logger log.DebugLogger) {
	err := amipublisher.ExpireResources(targets, skipTargets, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error expiring resources: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
