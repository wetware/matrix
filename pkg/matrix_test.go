package mx_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/config"
	"github.com/stretchr/testify/require"

	mock_mx "github.com/wetware/matrix/internal/mock/pkg"
	"github.com/wetware/matrix/internal/testutil"
	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/clock"
)

const n = 10

func TestSimulation(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim := mx.New(ctx,
		mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
		mx.WithHostFactory(testutil.NewHostFactory(ctrl)))

	t.Run("HostSet", func(t *testing.T) {
		t.Parallel()

		hs := sim.MustHostSet(ctx, n)
		require.Len(t, hs, 10)

		t.Run("EnsureUniquePeers", func(t *testing.T) {
			t.Parallel()

			// Ensure unique
			seen := make(map[host.Host]struct{})
			hs.Map(func(i int, h host.Host) error {
				if _, ok := seen[h]; ok {
					return errors.New("duplicate")
				}

				seen[h] = struct{}{}
				return nil
			})
		})

		t.Run("MustHostSetPanicsOnError", func(t *testing.T) {
			t.Parallel()

			hf := mock_mx.NewMockHostFactory(ctrl)
			hf.EXPECT().
				NewHost(
					gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
					gomock.AssignableToTypeOf([]config.Option(nil)),
				).
				DoAndReturn(func(ctx context.Context, opt []config.Option) (host.Host, error) {
					return nil, errors.New("test")
				}).
				Times(n) // sim.MustHostSet calls mx.Go under the hood

			sim := mx.New(ctx,
				mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
				mx.WithHostFactory(hf))

			require.Panics(t, func() {
				sim.MustHostSet(ctx, n)
			})
		})
	})

	t.Run("Clock", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, sim.Clock().Accuracy(), clock.DefaultAccuracy)
	})

}
