package matrix_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	matrix "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/discover"
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

		env.Op(matrix.Announce(discover.SelectAll{}, ns)).
			Then(matrix.Discover(discover.SelectAll{}, ns)).
			// Then(matrix.Filter()).
			Call(ctx, h0, h1).
			Must()
	})
}
