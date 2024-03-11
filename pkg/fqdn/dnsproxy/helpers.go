// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package dnsproxy

import (
	"fmt"
	"net"

	"github.com/cilium/dns"
)

// lookupTargetDNSServer finds the intended DNS target server for a specific
// request (passed in via ServeDNS). The IP:port combination is
// returned.
func lookupTargetDNSServer(w dns.ResponseWriter) (serverIP net.IP, serverPort, serverProtocol uint16, addrStr string, err error) {
	switch addr := (w.LocalAddr()).(type) {
	case *net.UDPAddr:
		return addr.IP, uint16(addr.Port), 17, addr.String(), nil
	case *net.TCPAddr:
		return addr.IP, uint16(addr.Port), 6, addr.String(), nil
	default:
		return nil, 0, 0, addr.String(), fmt.Errorf("Cannot extract address information for type %T: %+v", addr, addr)
	}
}
