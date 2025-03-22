package main

import (
	"fmt"
	"log"

	"github.com/rousage/filestorage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of the TCPTransporter")
	return nil
}

func main() {
	opts := p2p.TCPTransporterOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransporter(opts)

	go func() error {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.Listen(); err != nil {
		log.Fatal(err)
	}

	select {}
}
