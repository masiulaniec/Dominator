package rpcd

import (
	"github.com/masiulaniec/Dominator/lib/filter"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/proto/sub"
)

func (t *rpcType) SetConfiguration(conn *srpc.Conn,
	request sub.SetConfigurationRequest,
	reply *sub.SetConfigurationResponse) error {
	if request.CpuPercent > 100 {
		request.CpuPercent = 100
	}
	if request.CpuPercent > 0 {
		t.scannerConfiguration.DefaultCpuPercent = request.CpuPercent
		t.scannerConfiguration.CpuLimiter.SetCpuPercent(request.CpuPercent)
	}
	if request.NetworkSpeedPercent > 0 {
		t.scannerConfiguration.NetworkReaderContext.SetSpeedPercent(
			request.NetworkSpeedPercent)
	}
	if request.ScanSpeedPercent > 0 {
		t.scannerConfiguration.FsScanContext.GetContext().SetSpeedPercent(
			request.ScanSpeedPercent)
	}
	newFilter, err := filter.New(request.ScanExclusionList)
	if err != nil {
		return err
	}
	t.scannerConfiguration.ScanFilter = newFilter
	t.logger.Printf("SetConfiguration()\n")
	return nil
}
