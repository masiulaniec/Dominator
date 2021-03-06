package urlutil

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/masiulaniec/Dominator/lib/fsutil"
	"github.com/masiulaniec/Dominator/lib/log"
)

func watchUrl(rawurl string, checkInterval time.Duration,
	logger log.Logger) (<-chan io.ReadCloser, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "file" {
		return fsutil.WatchFile(u.Path, logger), nil
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		ch := make(chan io.ReadCloser, 1)
		go watchUrlLoop(rawurl, checkInterval, ch, logger)
		return ch, nil
	}
	return nil, errors.New("unknown scheme: " + u.Scheme)
}

func watchUrlLoop(rawurl string, checkInterval time.Duration,
	ch chan<- io.ReadCloser, logger log.Logger) {
	for ; ; time.Sleep(checkInterval) {
		watchUrlOnce(rawurl, ch, logger)
		if checkInterval <= 0 {
			return
		}
	}
}

func watchUrlOnce(rawurl string, ch chan<- io.ReadCloser, logger log.Logger) {
	resp, err := http.Get(rawurl)
	if err != nil {
		logger.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Println(resp.Status)
		return
	}
	ch <- resp.Body
}
