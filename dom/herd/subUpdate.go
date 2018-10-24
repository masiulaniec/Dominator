package herd

import (
	"syscall"
	"time"

	"github.com/masiulaniec/Dominator/dom/lib"
	subproto "github.com/masiulaniec/Dominator/proto/sub"
)

// Returns (idle, missing), idle=true if no update needs to be performed.
func (sub *Sub) buildUpdateRequest(request *subproto.UpdateRequest) (
	bool, bool) {
	request.ImageName = sub.requiredImageName
	request.Triggers = sub.requiredImage.Triggers
	var rusageStart, rusageStop syscall.Rusage
	computeStartTime := time.Now()
	syscall.Getrusage(syscall.RUSAGE_SELF, &rusageStart)
	subObj := lib.Sub{
		Hostname:       sub.mdb.Hostname,
		FileSystem:     sub.fileSystem,
		ComputedInodes: sub.computedInodes,
		ObjectCache:    sub.objectCache}
	if lib.BuildUpdateRequest(subObj, sub.requiredImage, request, false, false,
		sub.herd.logger) {
		return false, true
	}
	syscall.Getrusage(syscall.RUSAGE_SELF, &rusageStop)
	computeTimeDistribution.Add(time.Since(computeStartTime))
	sub.lastComputeUpdateCpuDuration = time.Duration(
		rusageStop.Utime.Sec)*time.Second +
		time.Duration(rusageStop.Utime.Usec)*time.Microsecond -
		time.Duration(rusageStart.Utime.Sec)*time.Second -
		time.Duration(rusageStart.Utime.Usec)*time.Microsecond
	computeCpuTimeDistribution.Add(sub.lastComputeUpdateCpuDuration)
	if len(request.FilesToCopyToCache) > 0 ||
		len(request.InodesToMake) > 0 ||
		len(request.HardlinksToMake) > 0 ||
		len(request.PathsToDelete) > 0 ||
		len(request.DirectoriesToMake) > 0 ||
		len(request.InodesToChange) > 0 ||
		sub.lastSuccessfulImageName != sub.requiredImageName {
		sub.herd.logger.Debugf(0,
			"buildUpdateRequest(%s) took: %s user CPU time\n",
			sub, sub.lastComputeUpdateCpuDuration)
		return false, false
	}
	return true, false
}
