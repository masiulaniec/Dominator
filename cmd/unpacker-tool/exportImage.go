package main

import (
	"fmt"
	"os"

	"github.com/masiulaniec/Dominator/imageunpacker/client"
	"github.com/masiulaniec/Dominator/lib/srpc"
)

func exportImageSubcommand(srpcClient *srpc.Client, args []string) {
	err := client.ExportImage(srpcClient, args[0], args[1], args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting image: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
