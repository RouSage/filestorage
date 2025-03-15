package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransporter(t *testing.T) {
	opts := TCPTransporterOpts{
		ListenAddr:    ":4000",
		HandshakeFunc: NOOPHandshakeFunc,
	}
	tr := NewTCPTransporter(opts)

	assert.Equal(t, tr.ListenAddr, opts.ListenAddr)
	assert.Nil(t, tr.Listen())
}
