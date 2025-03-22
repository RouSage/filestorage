package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpeg bytes"))
	if err := s.writeStream("mypicture", data); err != nil {
		t.Error(err)
	}
}
