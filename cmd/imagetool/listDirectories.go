package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imageserver/client"
	"github.com/masiulaniec/Dominator/lib/image"
	"github.com/masiulaniec/Dominator/lib/srpc"
)

func listDirectoriesSubcommand(args []string) {
	imageClient, _ := getClients()
	if err := listDirectories(imageClient); err != nil {
		fmt.Fprintf(os.Stderr, "Error listing directories: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func listDirectories(imageSClient *srpc.Client) error {
	directories, err := client.ListDirectories(imageSClient)
	if err != nil {
		return err
	}
	image.SortDirectories(directories)
	maxDirnameWidth := 0
	for _, directory := range directories {
		if len(directory.Name) > maxDirnameWidth {
			maxDirnameWidth = len(directory.Name)
		}
	}
	for _, directory := range directories {
		if directory.Metadata.OwnerGroup == "" {
			fmt.Println(directory.Name)
			continue
		}
		fmt.Printf("%-*s  ", maxDirnameWidth, directory.Name)
		fmt.Printf("OwnerGroup=%s", directory.Metadata.OwnerGroup)
		fmt.Println()
	}
	return nil
}
