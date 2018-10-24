package rpcd

import (
	"io"

	"github.com/masiulaniec/Dominator/fleetmanager/hypervisors"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/lib/srpc/serverutil"
)

type srpcType struct {
	hypervisorsManager *hypervisors.Manager
	logger             log.DebugLogger
	*serverutil.PerUserMethodLimiter
}

type htmlWriter srpcType

func (hw *htmlWriter) WriteHtml(writer io.Writer) {
	hw.writeHtml(writer)
}

func Setup(hypervisorsManager *hypervisors.Manager, logger log.DebugLogger) (
	*htmlWriter, error) {
	srpcObj := &srpcType{
		hypervisorsManager: hypervisorsManager,
		logger:             logger,
		PerUserMethodLimiter: serverutil.NewPerUserMethodLimiter(
			map[string]uint{
				"GetMachineInfo": 1,
				"GetUpdates":     1,
			}),
	}
	srpc.RegisterNameWithOptions("FleetManager", srpcObj,
		srpc.ReceiverOptions{
			PublicMethods: []string{
				"GetHypervisorForVM",
				"GetMachineInfo",
				"GetUpdates",
				"ListHypervisorLocations",
				"ListHypervisorsInLocation",
				"ListVMsInLocation",
			}})
	return (*htmlWriter)(srpcObj), nil
}
