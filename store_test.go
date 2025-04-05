package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "mybestpicture"
	data := []byte("some jpeg bytes")

	// write data to file
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)

	// read data from file
	r, err := s.Read(key)
	assert.Nil(t, err)

	b, err := io.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, data, b)
}

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathname := CASPathTransformFunc(key)

	expected := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	expectedOriginal := "be17b32c2870b1c0c73b59949db6a3be7814dd23"

	assert.Equal(t, expected, pathname.Pathname)
	assert.Equal(t, expectedOriginal, pathname.Filename)
}
