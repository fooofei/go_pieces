package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fooofei/sshttp"
	"github.com/gorilla/websocket"
)

func usage(program string) string {
	fmt1 := `
Usage:
	ssh -o "ProxyCommand %v proxy_server_host proxy_server_port %%h %%p" user@host
	proxy_server_host: hostname on which Proxy Server runs on
	proxy_server_port: TCP port number to connect to Proxy Server
	user: SSH user
	host: SSH host

Example:
	ssh -o "ProxyCommand %v 127.0.0.1 8888 %%h %%p" work@192.168.200.128
`
	return fmt.Sprintf(fmt1, program, program)
}

func main() {
	pro := os.Args[0]
	fmt.Printf("%v", usage(pro))

	host := os.Args[1]
	port := os.Args[2]
	sshdHost := os.Args[3]
	sshdPort := os.Args[4]

	addr := net.JoinHostPort(host, port)

	waitCtx, cancel := context.WithCancel(context.Background())
	d := net.Dialer{}
	conn, err := d.DialContext(waitCtx, "tcp", addr)
	if err != nil {
		log.Fatalf("dial err= %v", err)
	}
	websocket.Dialer{
		NetDial:           nil,
		NetDialContext:    nil,
		Proxy:             nil,
		TLSClientConfig:   nil,
		HandshakeTimeout:  0,
		ReadBufferSize:    0,
		WriteBufferSize:   0,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		EnableCompression: false,
		Jar:               nil,
	}()

	// first login
	_, _ = conn.Write(sshttp.NewLogin())

	_, _ = conn.Write(sshttp.NewProxyConnect(net.JoinHostPort(sshdHost, sshdPort)))

	_ = cancel
}
