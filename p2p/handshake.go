package p2p

type HandshakeFunc func(Peer) error

func NOOPHandshakeFunc(Peer) error { return nil }
