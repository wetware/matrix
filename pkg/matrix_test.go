package mx_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wetware/matrix/internal/testutil"
	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/netsim"
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

		sim.Op(mx.Announce(netsim.SelectAll{}, ns)).
			Then(mx.Discover(netsim.SelectAll{}, ns)).
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

func ExampleSimulation() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim := mx.New(ctx)

	h0 := sim.MustHost(ctx)
	h1 := sim.MustHost(ctx)

	h0.SetStreamHandler(echo, func(s network.Stream) {
		defer s.Close()

		if _, err := io.Copy(s, s); err != nil {
			panic(err)
		}
	})

	sim.Op(mx.Announce(netsim.SelectAll{}, ns)).
		Then(mx.Discover(netsim.SelectAll{}, ns)).
		Call(ctx, h0, h1).
		Must()

	s, err := h1.NewStream(ctx, h0.ID(), echo)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	_, err = io.Copy(s, strings.NewReader("Hello, world!"))
	if err != nil {
		panic(err)
	}
	s.CloseWrite()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, s)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
	// Output: Hello, world!
}

func TestNewDiscovery(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim := mx.New(ctx)

	h := testutil.NewHost(ctrl)

	svc := sim.NewDiscovery(h, nil)
	require.NotNil(t, svc)
	require.NotNil(t, svc.Topo)
	require.Equal(t, svc.Info, host.InfoFromHost(h))
}
