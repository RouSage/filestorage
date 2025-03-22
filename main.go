package main

import (
	"log"

	"github.com/rousage/filestorage/p2p"
)

func main() {
	opts := p2p.TCPTransporterOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransporter(opts)

	if err := tr.Listen(); err != nil {
		log.Fatal(err)
	}

	select {}
}
