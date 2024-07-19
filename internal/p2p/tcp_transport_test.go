package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":3000"
	tr := NewTCPTransport(listenAddr)
	assert.Equal(t, listenAddr, tr.ListenAddr())
}
