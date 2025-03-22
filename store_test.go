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

	assert.Equal(t, "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23", pathname)
}
