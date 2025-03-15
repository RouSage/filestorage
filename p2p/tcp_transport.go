package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransporter struct {
	listenAddr string
	listener   net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransporter(listenAddr string) *TCPTransporter {
	return &TCPTransporter{
		listenAddr: listenAddr,
	}
}

func (t *TCPTransporter) Listen() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransporter) startAcceptLoop() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}
}

func (t *TCPTransporter) handleConn(conn net.Conn) {
	fmt.Printf("new incoming connection %+v\n", conn)
}
