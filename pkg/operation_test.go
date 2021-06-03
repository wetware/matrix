package mx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/internal/testutil"
	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/netsim"
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

		res, err := mx.Select(func(ctx context.Context, hs mx.Selection) (mx.Selection, error) {
			called = true
			return hs, nil
		}).EvalArgs(ctx, hs...) // use 'EvalArgs' to maximize test coverage

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

			res, err := mx.Select(mx.Nop()).
				Bind(func(f mx.SelectFunc) mx.Op {
					return f.Bind(func(ctx context.Context, got mx.Selection) (mx.Selection, error) {
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
				res mx.Selection
				err error
			)

			require.NotPanics(t, func() {
				res, err = mx.Select(mx.Just(hs)).
					Bind(func(f mx.SelectFunc) mx.Op {
						return f.Bind(func(context.Context, mx.Selection) (mx.Selection, error) {
							return nil, errors.New("test")
						})
					}).
					Bind(func(f mx.SelectFunc) mx.Op {
						return f.Bind(func(context.Context, mx.Selection) (mx.Selection, error) {
							panic("not aborted")
						})
					}).
					Eval(ctx, hs)
			})

			require.Nil(t, res)
			require.EqualError(t, err, "test")
		})
	})

	t.Run("Map", func(t *testing.T) {
		t.Parallel()

		res := mx.Select(mx.Nop()).Map(func(_ context.Context, _ int, h host.Host) error {
			return nil
		}).Must(ctx, hs)
		require.ElementsMatch(t, hs, res)
	})

	t.Run("MustPanicsOnError", func(t *testing.T) {
		t.Parallel()

		require.Panics(t, func() {
			mx.Select(mx.Fail(errors.New("err"))).Must(ctx, hs)
		})
	})

	t.Run("Filter", func(t *testing.T) {
		t.Parallel()

		even := mx.Filter(func(i int, _ host.Host) bool {
			return i%2 == 0 // select even
		}).Must(ctx, hs)

		require.Len(t, even, len(hs)/2)
	})

	t.Run("DiscoverReturnsTopologyErrors", func(t *testing.T) {
		t.Parallel()

		sim := mx.New(ctx,
			mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
			mx.WithHostFactory(testutil.NewHostFactory(ctrl)))

		res, err := mx.Map(mx.Discover(sim, errTopology{}, "")).
			EvalArgs(ctx, hs[0])
		require.EqualError(t, err, "test")
		require.Nil(t, res)
	})
}

func mkHostSlice(ctx context.Context, ctrl *gomock.Controller) mx.Selection {
	return mx.New(ctx,
		mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
		mx.WithHostFactory(testutil.NewHostFactory(ctrl))).
		MustHostSet(ctx, n)

}

type errTopology struct{}

func (errTopology) SetDefaultOptions(*discovery.Options) error {
	return errors.New("test")
}

func (errTopology) Select(context.Context, netsim.Scope, *peer.AddrInfo, *discovery.Options) (netsim.InfoSlice, error) {
	panic("unreachable")
}
