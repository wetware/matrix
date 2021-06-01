package matrix_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	matrix "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/net"
)

const ns = "matrix.test"

func TestIntegration(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		env := matrix.New(ctx)

		h0 := env.MustHost(ctx)
		h1 := env.MustHost(ctx)

		env.Op(matrix.Announce(net.SelectAll{}, ns)).
			Then(matrix.Discover(net.SelectAll{}, ns)).
			// Then(matrix.Filter()).
			Call(ctx, h0, h1).
			Must()
	})
}
