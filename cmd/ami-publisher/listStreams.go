package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imagepublishers/amipublisher"
	libjson "github.com/masiulaniec/Dominator/lib/json"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/verstr"
)

func listStreamsSubcommand(args []string, logger log.DebugLogger) {
	err := listStreams(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing streams: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func listStreams(logger log.DebugLogger) error {
	results, err := amipublisher.ListStreams(targets, skipTargets,
		*instanceName, logger)
	if err != nil {
		return err
	}
	streams := make([]string, 0, len(results))
	for stream := range results {
		streams = append(streams, stream)
	}
	verstr.Sort(streams)
	return libjson.WriteWithIndent(os.Stdout, "    ", streams)
}
