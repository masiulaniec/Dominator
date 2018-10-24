package main

import (
	"github.com/masiulaniec/Dominator/lib/net/proxy"
	"github.com/masiulaniec/Dominator/lib/srpc"
)

func dialFleetManager(address string) (*srpc.Client, error) {
	return srpc.DialHTTP("tcp", address, 0)
}

func dialHypervisor(address string) (*srpc.Client, error) {
	if dialer, err := proxy.NewDialer(*hypervisorProxy); err != nil {
		return nil, err
	} else {
		return srpc.DialHTTPWithDialer("tcp", address, dialer)
	}
}
