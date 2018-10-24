package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imagepublishers/amipublisher"
	libjson "github.com/masiulaniec/Dominator/lib/json"
	"github.com/masiulaniec/Dominator/lib/log"
)

func listImagesSubcommand(args []string, logger log.DebugLogger) {
	err := listImages(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing images: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func listImages(logger log.DebugLogger) error {
	results, err := amipublisher.ListImages(targets, skipTargets, searchTags,
		excludeSearchTags, *minImageAge, logger)
	if err != nil {
		return err
	}
	return libjson.WriteWithIndent(os.Stdout, "    ", results)
}
