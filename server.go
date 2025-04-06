package main

import (
	"fmt"
	"io"
	"log"

	"github.com/rousage/filestorage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transporter       p2p.Transporter
}

type FileServer struct {
	FileServerOpts

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
	}
}

func (s *FileServer) Start() error {
	if err := s.Transporter.Listen(); err != nil {
		return err
	}

	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Store(key string, r io.Reader) error {
	return s.store.Write(key, r)
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
