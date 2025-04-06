package main

import "github.com/rousage/filestorage/p2p"

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transporter       p2p.Transporter
}

type FileServer struct {
	FileServerOpts

	store *Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transporter.Listen(); err != nil {
		return err
	}

	return nil
}
