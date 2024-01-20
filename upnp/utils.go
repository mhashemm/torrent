package upnp

import (
	"net"
	"net/netip"
)

func udpRequest(addr string, port int, payload []byte) ([]byte, error) {
	socket, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	defer socket.Close()
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	remote := &net.UDPAddr{
		IP:   ip.AsSlice(),
		Port: port,
	}
	_, err = socket.WriteToUDP(payload, remote)
	if err != nil {
		return nil, err
	}
	received := make([]byte, 4096)
	n, err := socket.Read(received)
	if err != nil {
		return nil, err
	}
	return received[:n], nil
}
