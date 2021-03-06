package main

import (
	"fmt"
	"net"
	"os"

	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/log"
	proto "github.com/masiulaniec/Dominator/proto/hypervisor"
)

func discardVmOldImageSubcommand(args []string, logger log.DebugLogger) {
	if err := discardVmOldImage(args[0], logger); err != nil {
		fmt.Fprintf(os.Stderr, "Error discarding VM old image: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func discardVmOldImage(vmHostname string, logger log.DebugLogger) error {
	if vmIP, hypervisor, err := lookupVmAndHypervisor(vmHostname); err != nil {
		return err
	} else {
		return discardVmOldImageOnHypervisor(hypervisor, vmIP, logger)
	}
}

func discardVmOldImageOnHypervisor(hypervisor string, ipAddr net.IP,
	logger log.DebugLogger) error {
	request := proto.DiscardVmOldImageRequest{ipAddr}
	client, err := dialHypervisor(hypervisor)
	if err != nil {
		return err
	}
	defer client.Close()
	var reply proto.DiscardVmOldImageResponse
	err = client.RequestReply("Hypervisor.DiscardVmOldImage", request, &reply)
	if err != nil {
		return err
	}
	return errors.New(reply.Error)
}
