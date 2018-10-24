package httpd

import (
	"bufio"
	"net/http"

	"github.com/masiulaniec/Dominator/lib/json"
)

func (s state) listAvailableAddressesHandler(w http.ResponseWriter,
	req *http.Request) {
	writer := bufio.NewWriter(w)
	defer writer.Flush()
	addresses := s.manager.ListAvailableAddresses()
	json.WriteWithIndent(writer, "    ", addresses)
}
