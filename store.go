package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
)

func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hasStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hasStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hasStr[from:to]
	}

	return path.Join(paths...)
}

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathname := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"
	fullPath := path.Join(pathname, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return nil
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullPath)

	return nil
}
