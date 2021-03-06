package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/masiulaniec/Dominator/lib/errors"
	"github.com/masiulaniec/Dominator/lib/json"
	"github.com/masiulaniec/Dominator/lib/srpc"
	//"github.com/masiulaniec/Dominator/lib/fsutil"
	"github.com/masiulaniec/Dominator/lib/log"
	proto "github.com/masiulaniec/Dominator/proto/hypervisor"
)

type devicesInfo struct {
	Interfaces []interfaceType `xml:"interface"`
	Volumes    []volumeType    `xml:"disk"`
}

type driverType struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type interfaceType struct {
	Mac  macType `xml:"mac"`
	Type string  `xml:"type,attr"`
}

type macType struct {
	Address string `xml:"address,attr"`
}

type memoryInfo struct {
	Value uint64 `xml:",chardata"`
	Unit  string `xml:"unit,attr"`
}

type sourceType struct {
	File string `xml:"file,attr"`
}

type vCpuInfo struct {
	Num       uint   `xml:",chardata"`
	Placement string `xml:"placement,attr"`
}

type virshInfoType struct {
	Devices devicesInfo `xml:"devices"`
	Memory  memoryInfo  `xml:"memory"`
	Name    string      `xml:"name"`
	VCpu    vCpuInfo    `xml:"vcpu"`
}

type volumeType struct {
	Device string     `xml:"device,attr"`
	Driver driverType `xml:"driver"`
	Source sourceType `xml:"source"`
	Type   string     `xml:"type,attr"`
}

func importVirshVmSubcommand(args []string, logger log.DebugLogger) {
	if err := importVirshVm(args[0], args[1], logger); err != nil {
		fmt.Fprintf(os.Stderr, "Error importing VM: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func ensureDomainIsStopped(domainName string) error {
	state, err := getDomainState(domainName)
	if err != nil {
		return err
	}
	if state == "shut off" {
		return nil
	}
	if state != "running" {
		return fmt.Errorf("domain is in unsupported state \"%s\"", state)
	}
	response, err := askForInputChoice("Cannot import running VM",
		[]string{"shutdown", "quit"})
	if response == "quit" {
		return fmt.Errorf("domain must be shut off but is \"%s\"", state)
	}
	err = exec.Command("virsh", []string{"shutdown", domainName}...).Run()
	if err != nil {
		return err
	}
	for ; ; time.Sleep(time.Second) {
		state, err := getDomainState(domainName)
		if err != nil {
			return err
		}
		if state == "shut off" {
			return nil
		}
	}
}

func getDomainState(domainName string) (string, error) {
	cmd := exec.Command("virsh", []string{"domstate", domainName}...)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

func importVirshVm(macAddr, domainName string, logger log.DebugLogger) error {
	ipList, err := net.LookupIP(domainName)
	if err != nil {
		return err
	}
	if len(ipList) != 1 {
		return fmt.Errorf("number of IPs %d != 1", len(ipList))
	}
	tags := vmTags.Copy()
	if _, ok := tags["Name"]; !ok {
		tags["Name"] = domainName
	}
	request := proto.ImportLocalVmRequest{VmInfo: proto.VmInfo{
		Hostname:    domainName,
		OwnerGroups: ownerGroups,
		OwnerUsers:  ownerUsers,
		Tags:        tags,
	}}
	request.VerificationCookie, err = readImportCookie(logger)
	if err != nil {
		return err
	}
	hypervisor := fmt.Sprintf(":%d", *hypervisorPortNum)
	client, err := srpc.DialHTTP("tcp", hypervisor, 0)
	if err != nil {
		return err
	}
	defer client.Close()
	directories, err := listVolumeDirectories(client)
	if err != nil {
		return err
	}
	volumeRoots := make(map[string]string, len(directories))
	for _, dirname := range directories {
		volumeRoots[filepath.Dir(dirname)] = dirname
	}
	if err := ensureDomainIsStopped(domainName); err != nil {
		return err
	}
	cmd := exec.Command("virsh",
		[]string{"dumpxml", "--inactive", domainName}...)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	var virshInfo virshInfoType
	if err := xml.Unmarshal(stdout, &virshInfo); err != nil {
		return err
	}
	if macAddr != virshInfo.Devices.Interfaces[0].Mac.Address {
		return fmt.Errorf("MAC address specified: %s != virsh data: %s",
			macAddr, virshInfo.Devices.Interfaces[0].Mac.Address)
	}
	json.WriteWithIndent(os.Stdout, "    ", virshInfo)
	if numIf := len(virshInfo.Devices.Interfaces); numIf != 1 {
		return fmt.Errorf("number of interfaces %d != 1", numIf)
	}
	request.VmInfo.Address = proto.Address{
		IpAddress:  ipList[0],
		MacAddress: virshInfo.Devices.Interfaces[0].Mac.Address,
	}
	switch virshInfo.Memory.Unit {
	case "KiB":
		request.VmInfo.MemoryInMiB = virshInfo.Memory.Value >> 10
	case "MiB":
		request.VmInfo.MemoryInMiB = virshInfo.Memory.Value
	case "GiB":
		request.VmInfo.MemoryInMiB = virshInfo.Memory.Value << 10
	default:
		return fmt.Errorf("unknown memory unit: %s", virshInfo.Memory.Unit)
	}
	request.VmInfo.MilliCPUs = virshInfo.VCpu.Num * 1000
	myPidStr := strconv.Itoa(os.Getpid())
	for index, inputVolume := range virshInfo.Devices.Volumes {
		if inputVolume.Device != "disk" {
			continue
		}
		var volumeFormat proto.VolumeFormat
		err := volumeFormat.UnmarshalText([]byte(inputVolume.Driver.Type))
		if err != nil {
			return err
		}
		inputFilename := inputVolume.Source.File
		var volumeRoot string
		for dirname := filepath.Dir(inputFilename); ; {
			if vr, ok := volumeRoots[dirname]; ok {
				volumeRoot = vr
				break
			}
			if dirname == "/" {
				break
			}
			dirname = filepath.Dir(dirname)
		}
		if volumeRoot == "" {
			return fmt.Errorf("no Hypervisor directory for: %s", inputFilename)
		}
		outputDirname := filepath.Join(volumeRoot, "import", myPidStr)
		if err := os.MkdirAll(outputDirname, dirPerms); err != nil {
			return err
		}
		defer os.RemoveAll(outputDirname)
		outputFilename := filepath.Join(outputDirname,
			fmt.Sprintf("volume-%d", index))
		if err := os.Link(inputFilename, outputFilename); err != nil {
			return err
		}
		request.VolumeFilenames = append(request.VolumeFilenames,
			outputFilename)
		request.VmInfo.Volumes = append(request.VmInfo.Volumes,
			proto.Volume{Format: volumeFormat})
	}
	requestWithoutSecrets := request
	requestWithoutSecrets.VerificationCookie = nil
	json.WriteWithIndent(os.Stdout, "    ", requestWithoutSecrets)
	var reply proto.GetVmInfoResponse
	err = client.RequestReply("Hypervisor.ImportLocalVm", request, &reply)
	if err != nil {
		return err
	}
	if err := errors.New(reply.Error); err != nil {
		return err
	}
	logger.Debugln(0, "imported VM")
	for _, dirname := range directories {
		os.RemoveAll(filepath.Join(dirname, "import", myPidStr))
	}
	if err := maybeWatchVm(client, hypervisor, ipList[0], logger); err != nil {
		return err
	}
	if err := askForCommitDecision(client, ipList[0]); err != nil {
		if err == errorCommitAbandoned {
			response, _ := askForInputChoice(
				"Do you want to restart the old VM", []string{"y", "n"})
			if response != "y" {
				return err
			}
			cmd = exec.Command("virsh", "start", domainName)
			if output, err := cmd.CombinedOutput(); err != nil {
				logger.Println(string(output))
				return err
			}
		}
		return err
	}
	defer virshInfo.deleteVolumes()
	cmd = exec.Command("virsh",
		[]string{"undefine", "--managed-save", "--snapshots-metadata",
			"--remove-all-storage", domainName}...)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Println(string(output))
		return err
	}
	return nil
}

func (virshInfo virshInfoType) deleteVolumes() {
	for _, inputVolume := range virshInfo.Devices.Volumes {
		if inputVolume.Device != "disk" {
			continue
		}
		os.Remove(inputVolume.Source.File)
	}
}
