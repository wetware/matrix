package mx_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/require"

	"github.com/wetware/matrix/internal/testutil"
	mx "github.com/wetware/matrix/pkg"
)

const n = 10

// func TestSimulation(t *testing.T) {
// 	t.Parallel()
// 	t.Helper()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	sim := mx.New(ctx,
// 		mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
// 		mx.WithHostFactory(testutil.NewHostFactory(ctrl)))

// 	t.Run("HostSet", func(t *testing.T) {
// 		t.Parallel()

// 		hs := sim.MustHostSet(ctx, n)
// 		require.Len(t, hs, 10)

// 		t.Run("EnsureUniquePeers", func(t *testing.T) {
// 			t.Parallel()

// 			// Ensure unique
// 			seen := make(map[host.Host]struct{})
// 			hs.Map(func(i int, h host.Host) error {
// 				if _, ok := seen[h]; ok {
// 					return errors.New("duplicate")
// 				}

// 				seen[h] = struct{}{}
// 				return nil
// 			})
// 		})

// 		t.Run("Select", func(t *testing.T) {
// 			sim.Op(mx.Select(func(ctx context.Context, _ mx.Simulation, hs mx.HostSlice) (mx.HostSlice, error) {
// 				return hs[:5], nil
// 			})).
// 				Then(mx.Select(func(ctx context.Context, sim mx.Simulation, new mx.HostSlice) (mx.HostSlice, error) {
// 					require.EqualValues(t, hs[:5], new)
// 					return new, nil
// 				})).
// 				Must(ctx, hs...)
// 		})
// 	})

// }

func TestNewDiscovery(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim := mx.New(ctx,
		mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
		mx.WithHostFactory(testutil.NewHostFactory(ctrl)))

	h := sim.MustHost(ctx)

	svc := sim.NewDiscovery(h, nil)
	require.NotNil(t, svc)
	require.NotNil(t, svc.Topo)
	require.Equal(t, svc.Info, host.InfoFromHost(h))
}
