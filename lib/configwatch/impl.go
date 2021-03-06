package configwatch

import (
	"bytes"
	"io"
	"time"

	"github.com/masiulaniec/Dominator/lib/fsutil"
	"github.com/masiulaniec/Dominator/lib/log"
	"github.com/masiulaniec/Dominator/lib/url/urlutil"
)

func watch(url string, checkInterval time.Duration,
	decoder Decoder, logger log.DebugLogger) (<-chan interface{}, error) {
	rawChannel, err := urlutil.WatchUrl(url, checkInterval, logger)
	if err != nil {
		return nil, err
	}
	configChannel := make(chan interface{}, 1)
	go watchLoop(rawChannel, configChannel, decoder, logger)
	return configChannel, nil
}

func watchLoop(rawChannel <-chan io.ReadCloser,
	configChannel chan<- interface{}, decoder Decoder, logger log.DebugLogger) {
	var previousChecksum []byte
	for reader := range rawChannel {
		checksumReader := fsutil.NewChecksumReader(reader)
		if config, err := decoder(checksumReader); err != nil {
			logger.Println(err)
		} else {
			newChecksum := checksumReader.GetChecksum()
			if bytes.Equal(newChecksum, previousChecksum) {
				logger.Debugln(1, "ignoring unchanged configuration")
			} else {
				configChannel <- config
				previousChecksum = newChecksum
			}
		}
		reader.Close()
	}
	close(configChannel)
}

func watchWithCache(url string, checkInterval time.Duration,
	decoder Decoder, cacheFilename string, initialTimeout time.Duration,
	logger log.DebugLogger) (<-chan interface{}, error) {
	rawChannel, err := urlutil.WatchUrlWithCache(url, checkInterval,
		cacheFilename, initialTimeout, logger)
	if err != nil {
		return nil, err
	}
	configChannel := make(chan interface{}, 1)
	go watchLoopWithCache(rawChannel, configChannel, decoder, logger)
	return configChannel, nil
}

func watchLoopWithCache(rawChannel <-chan *urlutil.CachedReadCloser,
	configChannel chan<- interface{}, decoder Decoder, logger log.DebugLogger) {
	var previousChecksum []byte
	for reader := range rawChannel {
		checksumReader := fsutil.NewChecksumReader(reader)
		if config, err := decoder(checksumReader); err != nil {
			logger.Println(err)
		} else {
			newChecksum := checksumReader.GetChecksum()
			if bytes.Equal(newChecksum, previousChecksum) {
				logger.Debugln(1, "ignoring unchanged configuration")
			} else {
				if err := reader.SaveCache(); err != nil {
					logger.Println(err)
				}
				configChannel <- config
				previousChecksum = newChecksum
			}
		}
		reader.Close()
	}
	close(configChannel)
}
