package p2p

import "net"

// Peer is an interface that represents the remote node
type Peer interface {
	RemoteAddr() net.Addr
	Close() error
}

// Transporter is anything that handles the communication
// between  the nodes in the network.
// This can be of the form (TCP, UDP, websockets, etc.)
type Transporter interface {
	Dial(string) error
	Listen() error
	Consume() <-chan RPC
	Close() error
}
