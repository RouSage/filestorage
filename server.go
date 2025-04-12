package main

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/rousage/filestorage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transporter       p2p.Transporter
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Store(key string, r io.Reader) error {
	return s.store.Write(key, r)
}

func (s *FileServer) Start() error {
	if err := s.Transporter.Listen(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s\n", p.RemoteAddr())

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			log.Printf("trying to connect with remote: %s\n", addr)
			if err := s.Transporter.Dial(addr); err != nil {
				log.Printf("dial error: %s\n", err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transporter.Close()
	}()

	for {
		select {
		case msg := <-s.Transporter.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}
