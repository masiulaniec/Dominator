package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imagepublishers/amipublisher"
	libjson "github.com/masiulaniec/Dominator/lib/json"
	"github.com/masiulaniec/Dominator/lib/log"
)

func startInstancesSubcommand(args []string, logger log.DebugLogger) {
	if err := startInstances(*instanceName, logger); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting instances: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func startInstances(name string, logger log.DebugLogger) error {
	results, err := amipublisher.StartInstances(targets, skipTargets, name,
		logger)
	if err != nil {
		return err
	}
	if err := libjson.WriteWithIndent(os.Stdout, "    ", results); err != nil {
		return err
	}
	for _, result := range results {
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
