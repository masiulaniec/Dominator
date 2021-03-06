package rpcd

import (
	"io"
	"sync"

	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/rateio"
	"github.com/masiulaniec/Dominator/lib/srpc"
	"github.com/masiulaniec/Dominator/sub/scanner"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
)

type rpcType struct {
	scannerConfiguration         *scanner.Configuration
	fileSystemHistory            *scanner.FileSystemHistory
	objectsDir                   string
	rootDir                      string
	networkReaderContext         *rateio.ReaderContext
	netbenchFilename             string
	oldTriggersFilename          string
	rescanObjectCacheFunction    func()
	disableScannerFunc           func(disableScanner bool)
	logger                       log.Logger
	rwLock                       sync.RWMutex
	getFilesLock                 sync.Mutex
	fetchInProgress              bool // Fetch() & Update() mutually exclusive.
	updateInProgress             bool
	startTimeNanoSeconds         int32 // For Fetch() or Update().
	startTimeSeconds             int64
	lastFetchError               error
	lastUpdateError              error
	lastUpdateHadTriggerFailures bool
	lastSuccessfulImageName      string
}

type addObjectsHandlerType struct {
	objectsDir           string
	scannerConfiguration *scanner.Configuration
	logger               log.Logger
}

type HtmlWriter struct {
	lastSuccessfulImageName *string
}

func Setup(configuration *scanner.Configuration, fsh *scanner.FileSystemHistory,
	objectsDirname string, rootDirname string,
	netReaderContext *rateio.ReaderContext,
	netbenchFname string, oldTriggersFname string,
	disableScannerFunction func(disableScanner bool),
	rescanObjectCacheFunction func(), logger log.Logger) *HtmlWriter {
	rpcObj := &rpcType{
		scannerConfiguration:      configuration,
		fileSystemHistory:         fsh,
		objectsDir:                objectsDirname,
		rootDir:                   rootDirname,
		networkReaderContext:      netReaderContext,
		netbenchFilename:          netbenchFname,
		oldTriggersFilename:       oldTriggersFname,
		rescanObjectCacheFunction: rescanObjectCacheFunction,
		disableScannerFunc:        disableScannerFunction,
		logger:                    logger}
	srpc.RegisterName("Subd", rpcObj)
	addObjectsHandler := &addObjectsHandlerType{
		objectsDir:           objectsDirname,
		scannerConfiguration: configuration,
		logger:               logger}
	srpc.RegisterName("ObjectServer", addObjectsHandler)
	tricorder.RegisterMetric("/image-name", &rpcObj.lastSuccessfulImageName,
		units.None, "name of the image for the last successful update")
	return &HtmlWriter{&rpcObj.lastSuccessfulImageName}
}

func (hw *HtmlWriter) WriteHtml(writer io.Writer) {
	hw.writeHtml(writer)
}
