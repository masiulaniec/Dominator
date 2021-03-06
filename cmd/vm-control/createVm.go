package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/masiulaniec/Dominator/lib/flagutil"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/srpc"
	fm_proto "github.com/masiulaniec/Dominator/proto/fleetmanager"
	hyper_proto "github.com/masiulaniec/Dominator/proto/hypervisor"
)

func init() {
	rand.Seed(time.Now().Unix() + time.Now().UnixNano())
}

func createVmSubcommand(args []string, logger log.DebugLogger) {
	if err := createVm(logger); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating VM: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func acknowledgeVm(client *srpc.Client, ipAddress net.IP) error {
	request := hyper_proto.AcknowledgeVmRequest{ipAddress}
	var reply hyper_proto.AcknowledgeVmResponse
	return client.RequestReply("Hypervisor.AcknowledgeVm", request, &reply)
}

func callCreateVm(client *srpc.Client, request hyper_proto.CreateVmRequest,
	reply *hyper_proto.CreateVmResponse, imageReader, userDataReader io.Reader,
	logger log.DebugLogger) error {
	conn, err := client.Call("Hypervisor.CreateVm")
	if err != nil {
		return fmt.Errorf("error calling Hypervisor.CreateVm: %s", err)
	}
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	if err := encoder.Encode(request); err != nil {
		return fmt.Errorf("error encoding request: %s", err)
	}
	// Stream any required data.
	if imageReader != nil {
		logger.Debugln(0, "uploading image")
		if _, err := io.Copy(conn, imageReader); err != nil {
			return fmt.Errorf("error uploading image: %s", err)
		}
	}
	if userDataReader != nil {
		logger.Debugln(0, "uploading user data")
		if _, err := io.Copy(conn, userDataReader); err != nil {
			return fmt.Errorf("error uploading user data: %s", err)
		}
	}
	if err := conn.Flush(); err != nil {
		return fmt.Errorf("error flushing: %s", err)
	}
	for {
		var response hyper_proto.CreateVmResponse
		if err := decoder.Decode(&response); err != nil {
			return fmt.Errorf("error decoding: %s", err)
		}
		if response.Error != "" {
			return errors.New(response.Error)
		}
		if response.ProgressMessage != "" {
			logger.Debugln(0, response.ProgressMessage)
		}
		if response.Final {
			*reply = response
			return nil
		}
	}
}

func createVm(logger log.DebugLogger) error {
	if hypervisor, err := getHypervisorAddress(); err != nil {
		return err
	} else {
		logger.Debugf(0, "creating VM on %s\n", hypervisor)
		return createVmOnHypervisor(hypervisor, logger)
	}
}

func createVmOnHypervisor(hypervisor string, logger log.DebugLogger) error {
	request := hyper_proto.CreateVmRequest{
		DhcpTimeout: *dhcpTimeout,
		VmInfo: hyper_proto.VmInfo{
			Hostname:    *vmHostname,
			MemoryInMiB: uint64(memory),
			MilliCPUs:   *milliCPUs,
			OwnerGroups: ownerGroups,
			OwnerUsers:  ownerUsers,
			Tags:        vmTags,
			SubnetId:    *subnetId,
		},
		MinimumFreeBytes: uint64(minFreeBytes),
		RoundupPower:     *roundupPower,
	}
	if sizes, err := parseSizes(secondaryVolumeSizes); err != nil {
		return err
	} else {
		request.SecondaryVolumes = sizes
	}
	var imageReader, userDataReader io.Reader
	if *imageName != "" {
		request.ImageName = *imageName
		request.ImageTimeout = *imageTimeout
	} else if *imageURL != "" {
		request.ImageURL = *imageURL
	} else if *imageFile != "" {
		file, size, err := getReader(*imageFile)
		if err != nil {
			return err
		} else {
			defer file.Close()
			request.ImageDataSize = uint64(size)
			imageReader = bufio.NewReader(io.LimitReader(file, size))
		}
	} else {
		return errors.New("no image specified")
	}
	if *userDataFile != "" {
		file, size, err := getReader(*userDataFile)
		if err != nil {
			return err
		} else {
			defer file.Close()
			request.UserDataSize = uint64(size)
			userDataReader = bufio.NewReader(io.LimitReader(file, size))
		}
	}
	client, err := dialHypervisor(hypervisor)
	if err != nil {
		return err
	}
	defer client.Close()
	var reply hyper_proto.CreateVmResponse
	err = callCreateVm(client, request, &reply, imageReader, userDataReader,
		logger)
	if err != nil {
		return err
	}
	if err := acknowledgeVm(client, reply.IpAddress); err != nil {
		return fmt.Errorf("error acknowledging VM: %s", err)
	}
	fmt.Println(reply.IpAddress)
	if reply.DhcpTimedOut {
		return errors.New("DHCP ACK timed out")
	}
	if *dhcpTimeout > 0 {
		logger.Debugln(0, "Received DHCP ACK")
	}
	return maybeWatchVm(client, hypervisor, reply.IpAddress, logger)
}

func getHypervisorAddress() (string, error) {
	if *hypervisorHostname != "" {
		return fmt.Sprintf("%s:%d", *hypervisorHostname, *hypervisorPortNum),
			nil
	}
	client, err := dialFleetManager(fmt.Sprintf("%s:%d",
		*fleetManagerHostname, *fleetManagerPortNum))
	if err != nil {
		return "", err
	}
	defer client.Close()
	if *adjacentVM != "" {
		if adjacentVmIpAddr, err := lookupIP(*adjacentVM); err != nil {
			return "", err
		} else {
			return findHypervisorClient(client, adjacentVmIpAddr)
		}
	}
	request := fm_proto.ListHypervisorsInLocationRequest{
		Location: *location,
		SubnetId: *subnetId,
	}
	var reply fm_proto.ListHypervisorsInLocationResponse
	err = client.RequestReply("FleetManager.ListHypervisorsInLocation",
		request, &reply)
	if err != nil {
		return "", err
	}
	if reply.Error != "" {
		return "", errors.New(reply.Error)
	}
	if numHyper := len(reply.HypervisorAddresses); numHyper < 1 {
		return "", errors.New("no active Hypervisors in location")
	} else if numHyper < 2 {
		return reply.HypervisorAddresses[0], nil
	} else {
		return reply.HypervisorAddresses[rand.Intn(numHyper-1)], nil
	}
}

func getReader(filename string) (io.ReadCloser, int64, error) {
	if file, err := os.Open(filename); err != nil {
		return nil, -1, err
	} else {
		fi, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, -1, err
		}
		return file, fi.Size(), nil
	}
}

func parseSizes(strSizes flagutil.StringList) ([]hyper_proto.Volume, error) {
	var volumes []hyper_proto.Volume
	for _, strSize := range strSizes {
		var size uint64
		if _, err := fmt.Sscanf(strSize, "%dM", &size); err == nil {
			volumes = append(volumes, hyper_proto.Volume{Size: size << 20})
		} else if _, err := fmt.Sscanf(strSize, "%dG", &size); err == nil {
			volumes = append(volumes, hyper_proto.Volume{Size: size << 30})
		} else {
			return nil,
				fmt.Errorf("error parsing secondary volume sizes: %s", err)
		}
	}
	return volumes, nil
}
