package main

import (
	"bytes"
	"log"
	"time"

	"github.com/rousage/filestorage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransporterOpts := p2p.TCPTransporterOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransporter := p2p.NewTCPTransporter(tcpTransporterOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transporter:       tcpTransporter,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransporter.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")

	go func() {
		if err := s1.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	data := bytes.NewReader([]byte("my big data file"))

	s2.Store("myprivatedata", data)

	select {}
}
