package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransporter(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransporter(listenAddr)

	assert.Equal(t, tr.listenAddr, listenAddr)
	assert.Nil(t, tr.Listen())
}
