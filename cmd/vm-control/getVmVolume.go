package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/rsync"
	proto "github.com/masiulaniec/Dominator/proto/hypervisor"
)

func getVmVolumeSubcommand(args []string, logger log.DebugLogger) {
	if err := getVmVolume(args[0], logger); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting VM volume: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func getVmVolume(vmHostname string, logger log.DebugLogger) error {
	if vmIP, hypervisor, err := lookupVmAndHypervisor(vmHostname); err != nil {
		return err
	} else {
		return getVmVolumeOnHypervisor(hypervisor, vmIP, logger)
	}
}

func getVmVolumeOnHypervisor(hypervisor string, ipAddr net.IP,
	logger log.DebugLogger) error {
	client, err := dialHypervisor(hypervisor)
	if err != nil {
		return err
	}
	defer client.Close()
	vmInfo, err := getVmInfoClient(client, ipAddr)
	if err != nil {
		return err
	}
	if *volumeIndex >= uint(len(vmInfo.Volumes)) {
		return fmt.Errorf("volumeIndex too large")
	}
	var initialFileSize uint64
	reader, err := os.OpenFile(*volumeFilename, os.O_RDONLY, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		defer reader.Close()
		if fi, err := reader.Stat(); err != nil {
			return err
		} else {
			initialFileSize = uint64(fi.Size())
			if initialFileSize > vmInfo.Volumes[*volumeIndex].Size {
				return errors.New("file larger than volume")
			}
		}
	}
	writer, err := os.OpenFile(*volumeFilename, os.O_WRONLY|os.O_CREATE,
		privateFilePerms)
	if err != nil {
		return err
	}
	defer writer.Close()
	request := proto.GetVmVolumeRequest{
		IpAddress:   ipAddr,
		VolumeIndex: *volumeIndex,
	}
	conn, err := client.Call("Hypervisor.GetVmVolume")
	if err != nil {
		if reader == nil {
			os.Remove(*volumeFilename)
		}
		return err
	}
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	if err := encoder.Encode(request); err != nil {
		return fmt.Errorf("error encoding request: %s", err)
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	var response proto.GetVmVolumeResponse
	if err := decoder.Decode(&response); err != nil {
		return err
	}
	if err := errors.New(response.Error); err != nil {
		return err
	}
	stats, err := rsync.GetBlocks(conn, decoder, encoder, reader, writer,
		vmInfo.Volumes[*volumeIndex].Size, initialFileSize)
	if err != nil {
		return err
	}
	logger.Debugf(0, "sent %d B, received %d/%d B (%.0f * speedup)\n",
		stats.NumWritten, stats.NumRead, vmInfo.Volumes[*volumeIndex].Size,
		float64(vmInfo.Volumes[*volumeIndex].Size)/
			float64(stats.NumRead+stats.NumWritten))
	return nil
}
