package main

import (
	"os"

	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
	"github.com/masiulaniec/Dominator/sub/client"
)

func setConfigSubcommand(getSubClient getSubClientFunc, args []string) {
	if err := setConfig(getSubClient()); err != nil {
		logger.Fatalf("Error setting config: %s\n", err)
	}
	os.Exit(0)
}

func setConfig(srpcClient *srpc.Client) error {
	var config sub.Configuration
	config.CpuPercent = *cpuPercent
	config.NetworkSpeedPercent = *networkSpeedPercent
	config.ScanExclusionList = scanExcludeList
	config.ScanSpeedPercent = *scanSpeedPercent
	return client.SetConfiguration(srpcClient, config)
}
