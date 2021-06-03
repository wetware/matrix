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
	"github.com/wetware/matrix/pkg/namespace"
	"github.com/wetware/matrix/pkg/netsim"
)

const ns = "matrix.test"

func TestDiscovery(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		c  = testutil.NewClock(ctrl, 0, nil)
		np = namespace.New(c)

		sim = mx.New(ctx,
			mx.WithClock(c),
			mx.WithHostFactory(testutil.NewHostFactory(ctrl)),
			mx.WithNamespaceFactory(func(mx.Clock) netsim.NamespaceProvider { return np }))

		h = sim.MustHost(ctx)
	)

	t.Run("NewDiscovery", func(t *testing.T) {
		t.Parallel()

		svc := sim.NewDiscovery(h, nil)
		require.NotNil(t, svc)
		require.NotNil(t, svc.Topo)
		require.Equal(t, svc.Info, host.InfoFromHost(h))
	})

	t.Run("TopologyError", func(t *testing.T) {
		res, err := mx.Map(mx.Discover(sim, errTopology{}, ns)).
			EvalArgs(ctx, h)
		require.EqualError(t, err, "test")
		require.Nil(t, res)
	})

	t.Run("NotAnnounced", func(t *testing.T) {
		res := mx.Map(mx.Discover(sim, netsim.SelectAll{}, ns)).
			MustArgs(ctx, h)
		require.ElementsMatch(t, mx.Selection{h}, res)

		s, ok := np.Load(ns)
		require.False(t, ok)
		require.Nil(t, s)
	})

	t.Run("Success", func(t *testing.T) {
		c.EXPECT().
			After(netsim.DefaultTTL, gomock.AssignableToTypeOf(func() {}))

		ttl := np.LoadOrCreate(ns).Upsert(host.InfoFromHost(h), &discovery.Options{Ttl: netsim.DefaultTTL})
		require.Equal(t, netsim.DefaultTTL, ttl)

		res := mx.Map(mx.Discover(sim, netsim.SelectAll{}, ns)).
			MustArgs(ctx, h)
		require.ElementsMatch(t, mx.Selection{h}, res)

		s, ok := np.Load(ns)
		require.True(t, ok)
		require.NotNil(t, s)

		require.ElementsMatch(t, netsim.InfoSlice{host.InfoFromHost(h)}, s.Peers())
	})
}

type errTopology struct{}

func (errTopology) SetDefaultOptions(*discovery.Options) error {
	return errors.New("test")
}

func (errTopology) Select(context.Context, netsim.Scope, *peer.AddrInfo, *discovery.Options) (netsim.InfoSlice, error) {
	panic("unreachable")
}
