package utils

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func HTTPRequestString(r http.Request) string {
	b := &strings.Builder{}
	b.WriteString(r.Method)
	b.WriteString(" * HTTP/1.1\r\n")
	for k, v := range r.Header {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, "; ")))
	}
	b.WriteString("\r\n")
	return b.String()
}

func LocalIPAddr() string {
	conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{6, 9, 6, 9},
		Port: 6969,
	})
	conn.SetDeadline(time.Unix(0, 0))
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
