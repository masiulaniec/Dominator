package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imageserver/client"
)

func checkImageSubcommand(args []string) {
	imageSClient, _ := getClients()
	imageExists, err := client.CheckImage(imageSClient, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking image\t%s\n", err)
		os.Exit(1)
	}
	if imageExists {
		os.Exit(0)
	}
	os.Exit(1)
}
