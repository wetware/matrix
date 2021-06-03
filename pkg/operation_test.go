package mx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/internal/testutil"
	mx "github.com/wetware/matrix/pkg"
)

func TestOperation(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hs := mkHostSlice(ctx, ctrl)

	t.Run("Eval", func(t *testing.T) {
		t.Parallel()

		var called bool

		res, err := mx.Op(func(ctx context.Context, hs mx.HostSlice) (mx.HostSlice, error) {
			called = true
			return hs, nil
		}).Eval(ctx, hs)

		require.NoError(t, err)
		require.True(t, called, "bound function not called")
		require.ElementsMatch(t, hs, res)
	})

	t.Run("Bind", func(t *testing.T) {
		t.Parallel()
		t.Helper()

		t.Run("Succeed", func(t *testing.T) {
			t.Parallel()

			var called bool

			res, err := mx.Op(mx.Nop()).
				Bind(func(f mx.OpFunc) mx.Operation {
					return f.Bind(func(ctx context.Context, got mx.HostSlice) (mx.HostSlice, error) {
						called = true
						require.Len(t, got, n)
						return hs, nil
					})
				}).
				Eval(ctx, hs)

			require.NoError(t, err)
			require.True(t, called, "bound function not called")
			require.ElementsMatch(t, hs, res)

		})

		t.Run("FailureAborts", func(t *testing.T) {
			t.Parallel()

			var (
				res mx.HostSlice
				err error
			)

			require.NotPanics(t, func() {
				res, err = mx.Op(mx.Just(hs)).
					Bind(func(f mx.OpFunc) mx.Operation {
						return f.Bind(func(context.Context, mx.HostSlice) (mx.HostSlice, error) {
							return nil, errors.New("test")
						})
					}).
					Bind(func(f mx.OpFunc) mx.Operation {
						return f.Bind(func(context.Context, mx.HostSlice) (mx.HostSlice, error) {
							panic("not aborted")
						})
					}).
					Eval(ctx, hs)
			})

			require.Nil(t, res)
			require.EqualError(t, err, "test")
		})
	})
}

func mkHostSlice(ctx context.Context, ctrl *gomock.Controller) mx.HostSlice {
	return mx.New(ctx,
		mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
		mx.WithHostFactory(testutil.NewHostFactory(ctrl))).
		MustHostSet(ctx, n)

}
