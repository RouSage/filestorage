package main

import (
	"log"
	"time"

	"github.com/rousage/filestorage/p2p"
)

func main() {
	tcpTransporterOpts := p2p.TCPTransporterOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: OnPeer func
	}
	tcpTransporter := p2p.NewTCPTransporter(tcpTransporterOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transporter:       tcpTransporter,
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
