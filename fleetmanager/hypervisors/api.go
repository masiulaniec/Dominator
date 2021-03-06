package hypervisors

import (
	"io"
	"net"
	"sync"

	"github.com/masiulaniec/Dominator/fleetmanager/topology"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/lib/tags"
	fm_proto "github.com/masiulaniec/Dominator/proto/fleetmanager"
	hyper_proto "github.com/masiulaniec/Dominator/proto/hypervisor"
)

const (
	probeStatusNotYetProbed probeStatus = iota
	probeStatusConnected
	probeStatusNoSrpc
	probeStatusNoService
	probeStatusBad
)

type hypervisorType struct {
	logger          log.DebugLogger
	mutex           sync.RWMutex
	conn            *srpc.Conn
	deleteScheduled bool
	healthStatus    string
	localTags       tags.Tags
	location        string
	machine         *fm_proto.Machine
	migratingVms    map[string]*vmInfoType // Key: VM IP address.
	probeStatus     probeStatus
	subnets         []hyper_proto.Subnet
	vms             map[string]*vmInfoType // Key: VM IP address.
}

type ipStorer interface {
	AddIPsForHypervisor(hypervisor net.IP, addrs []net.IP) error
	CheckIpIsRegistered(addr net.IP) (bool, error)
	SetIPsForHypervisor(hypervisor net.IP, addrs []net.IP) error
	UnregisterHypervisor(hypervisor net.IP) error
}

type locationType struct {
	notifiers map[<-chan fm_proto.Update]chan<- fm_proto.Update
}

type Manager struct {
	storer       Storer
	logger       log.DebugLogger
	invertTable  [256]byte
	mutex        sync.RWMutex
	topology     *topology.Topology
	hypervisors  map[string]*hypervisorType // Key: hypervisor machine name.
	locations    map[string]*locationType   // Key: location.
	migratingIPs map[string]struct{}        // Key: VM IP address.
	notifiers    map[<-chan fm_proto.Update]*locationType
	subnets      map[string]*subnetType // Key: Gateway IP.
	vms          map[string]*vmInfoType // Key: VM IP address.
}

type probeStatus uint

type Storer interface {
	ipStorer
	tagsStorer
	vmStorer
}

type subnetType struct {
	subnet  *topology.Subnet
	startIp net.IP
	stopIp  net.IP
	nextIp  net.IP
}

type tagsStorer interface {
	ReadMachineTags(hypervisor net.IP) (tags.Tags, error)
	WriteMachineTags(hypervisor net.IP, tgs tags.Tags) error
}

type vmInfoType struct {
	ipAddr string
	hyper_proto.VmInfo
	hypervisor *hypervisorType
}

type vmStorer interface {
	DeleteVm(hypervisor net.IP, ipAddr string) error
	ListVMs(hypervisor net.IP) ([]string, error)
	ReadVm(hypervisor net.IP, ipAddr string) (*hyper_proto.VmInfo, error)
	WriteVm(hypervisor net.IP, ipAddr string, vmInfo hyper_proto.VmInfo) error
}

func New(storer Storer, logger log.DebugLogger) (*Manager, error) {
	return newManager(storer, logger)
}

func (m *Manager) ChangeMachineTags(hostname string, tgs tags.Tags) error {
	return m.changeMachineTags(hostname, tgs)
}

func (m *Manager) CloseUpdateChannel(channel <-chan fm_proto.Update) {
	m.closeUpdateChannel(channel)
}

func (m *Manager) GetHypervisorForVm(ipAddr net.IP) (string, error) {
	return m.getHypervisorForVm(ipAddr)
}

func (m *Manager) GetMachineInfo(hostname string) (fm_proto.Machine, error) {
	return m.getMachineInfo(hostname)
}

func (m *Manager) GetTopology() (*topology.Topology, error) {
	return m.getTopology()
}

func (m *Manager) ListHypervisorsInLocation(
	request fm_proto.ListHypervisorsInLocationRequest) ([]string, error) {
	return m.listHypervisorsInLocation(request)
}

func (m *Manager) ListLocations(dirname string) ([]string, error) {
	return m.listLocations(dirname)
}

func (m *Manager) ListVMsInLocation(dirname string) ([]net.IP, error) {
	return m.listVMsInLocation(dirname)
}

func (m *Manager) MakeUpdateChannel(location string) <-chan fm_proto.Update {
	return m.makeUpdateChannel(location)
}

func (m *Manager) WriteHtml(writer io.Writer) {
	m.writeHtml(writer)
}

func (m *Manager) UpdateTopology(t *topology.Topology) {
	m.updateTopology(t)
}
