package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP connection.
type TCPPeer struct {
	// The underlying connection of the peer, TCP connection in this case.
	net.Conn

	// if dial and retrieve a conn => outbound == true
	//
	// if accept and retrieve a conn => outbound == false
	outbound bool

	Wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

// Send implements the Peer interface, which will write the payload to the
// underlying connection.
func (p *TCPPeer) Send(payload []byte) error {
	_, err := p.Conn.Write(payload)
	return err
}

type TCPTransporterOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransporter struct {
	TCPTransporterOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransporter(opts TCPTransporterOpts) *TCPTransporter {
	return &TCPTransporter{
		TCPTransporterOpts: opts,
		rpcch:              make(chan RPC),
	}
}

// Consume implements the Transporter interface, which will return a read-only channel
// for reading the incoming messages from another peer in the network.
func (t *TCPTransporter) Consume() <-chan RPC {
	return t.rpcch
}

// Close implements the Transporter interface, which will close the underlying listener.
func (t *TCPTransporter) Close() error {
	return t.listener.Close()
}

// Dial implements the Transporter interface.
func (t *TCPTransporter) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransporter) Listen() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transporter is listening on port: %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransporter) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransporter) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			fmt.Printf("TCP read error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr().String()

		peer.Wg.Add(1)
		fmt.Println("waiting until stream is done")

		t.rpcch <- rpc

		peer.Wg.Wait()
		fmt.Println("stream is done, continuing normal read loop")
	}
}
