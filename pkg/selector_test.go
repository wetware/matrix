package mx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/require"
	mx "github.com/wetware/matrix/pkg"
)

func TestMap(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hs := mkHostSlice(ctx, ctrl)

	t.Run("Succeed", func(t *testing.T) {
		res := make(mx.Selection, len(hs))
		mx.Map(func(_ context.Context, i int, h host.Host) error {
			res[i] = h
			return nil
		}).Must(ctx, hs)

		require.ElementsMatch(t, hs, res)
	})

	t.Run("FailureAborts", func(t *testing.T) {
		var check mx.Selection
		res, err := mx.Map(func(_ context.Context, i int, h host.Host) error {
			check = append(check, h)
			return errors.New("test")
		}).Eval(ctx, hs)

		require.EqualError(t, err, "test")
		require.Nil(t, res)
		require.ElementsMatch(t, mx.Selection{hs[0]}, check)
	})
}

func TestFail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hs := mkHostSlice(ctx, ctrl)

	res, err := mx.Select(mx.Fail(errors.New("test"))).
		Eval(ctx, hs)
	require.EqualError(t, err, "test")
	require.Nil(t, res)
}

func TestPartition(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hs := mkHostSlice(ctx, ctrl)

	t.Run("Num", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < n; i++ {
			require.Equal(t, i, mx.Partition(i).Num())
		}
	})

	t.Run("Singletons", func(t *testing.T) {
		t.Parallel()

		p := mx.Partition(n)

		for i := 0; i < p.Num(); i++ {
			res := mx.Select(p.Get(i)).Must(ctx, hs)
			require.Len(t, res, 1)
			require.Equal(t, hs[i], res[0])
		}
	})

	t.Run("Bipartite", func(t *testing.T) {
		t.Parallel()

		p := mx.Partition(2)

		var ps []mx.Selection
		for i := 0; i < 2; i++ {
			res := mx.Select(p.Get(i)).
				Then(func(_ context.Context, hs mx.Selection) (mx.Selection, error) {
					ps = append(ps, hs)
					return hs, nil
				}).
				Must(ctx, hs)

			require.Len(t, res, n/2)
		}

		require.Len(t, ps, 2)

		// TODO:  check that the partitions don't overlap
	})
}
