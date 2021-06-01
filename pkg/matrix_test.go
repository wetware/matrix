package mx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/net"
)

const (
	ns   = "matrix.test"
	echo = "/echo"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sim := mx.New(ctx)

		h0 := sim.MustHost(ctx)
		h1 := sim.MustHost(ctx)

		sim.Op(mx.Announce(net.SelectAll{}, ns)).
			Then(mx.Discover(net.SelectAll{}, ns)).
			Call(ctx, h0, h1).
			Must()
	})
}
