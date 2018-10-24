package main

import (
	"os"

	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
	"github.com/masiulaniec/Dominator/sub/client"
)

func deleteSubcommand(getSubClient getSubClientFunc, args []string) {
	if err := deletePaths(getSubClient(), args); err != nil {
		logger.Fatalf("Error deleting: %s\n", err)
	}
	os.Exit(0)
}

func deletePaths(srpcClient *srpc.Client, pathnames []string) error {
	return client.CallUpdate(srpcClient, sub.UpdateRequest{
		PathsToDelete: pathnames,
		Wait:          true},
		&sub.UpdateResponse{})
}
