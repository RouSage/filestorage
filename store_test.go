package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpeg bytes"))
	err := s.writeStream("mypicture", data)

	assert.Nil(t, err)
}

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathname := CASPathTransformFunc(key)

	expected := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	expectedOriginal := "be17b32c2870b1c0c73b59949db6a3be7814dd23"

	assert.Equal(t, expected, pathname.Pathname)
	assert.Equal(t, expectedOriginal, pathname.Original)
}
