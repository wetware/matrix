package mx_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/net"
)

const (
	ns   = "matrix.test"
	echo = "/echo"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var h0, h1 host.Host
	assert.NotPanics(t, func() {

		sim := mx.New(ctx)

		h0 = sim.MustHost(ctx)
		h0.SetStreamHandler(echo, func(s network.Stream) {
			defer s.Close()

			_, err := io.Copy(s, s)
			require.NoError(t, err)
		})

		h1 = sim.MustHost(ctx)

		sim.Op(mx.Announce(net.SelectAll{}, ns)).
			Then(mx.Discover(net.SelectAll{}, ns)).
			Call(ctx, h0, h1).
			Must()
	})

	s, err := h1.NewStream(ctx, h0.ID(), echo)
	require.NoError(t, err)

	var buf bytes.Buffer
	_, err = io.Copy(s, strings.NewReader("Hello, world!"))
	require.NoError(t, err)
	require.NoError(t, s.CloseWrite())

	_, err = io.Copy(&buf, s)
	require.NoError(t, err)

	assert.Equal(t, "Hello, world!", buf.String())
	assert.NoError(t, s.Close())
}
