package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if dial and retrieve a conn => outbound == true
	//
	// if accept and retrieve a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransporterOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransporter struct {
	TCPTransporterOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransporter(opts TCPTransporterOpts) *TCPTransporter {
	return &TCPTransporter{
		TCPTransporterOpts: opts,
	}
}

func (t *TCPTransporter) Listen() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
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

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransporter) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP decode error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()

		fmt.Printf("received message: %+v\n", rpc)
	}
}
