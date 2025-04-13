package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo-%d", i)
		data := []byte("some jpeg bytes")

		// write data to file
		n, err := s.Write(key, bytes.NewReader(data))
		assert.Nil(t, err)
		assert.Equal(t, int64(len(data)), n)

		// file should exist
		exists := s.Has(key)
		assert.Equal(t, true, exists)

		// read data from file
		r, err := s.Read(key)
		assert.Nil(t, err)

		b, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, data, b)

		// delete file
		err = s.Delete(key)
		assert.Nil(t, err)

		// file should not exist
		exists = s.Has(key)
		assert.Equal(t, false, exists)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathname := CASPathTransformFunc(key)

	extepctedPathname := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	expectedFilename := "be17b32c2870b1c0c73b59949db6a3be7814dd23"

	assert.Equal(t, extepctedPathname, pathname.Pathname)
	assert.Equal(t, expectedFilename, pathname.Filename)
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	err := s.Clear()
	assert.Nil(t, err)
}
