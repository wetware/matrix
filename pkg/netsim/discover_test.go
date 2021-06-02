package netsim_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/internal/testutil"
	"github.com/wetware/matrix/pkg/namespace"
	"github.com/wetware/matrix/pkg/netsim"
)

const n = 10

func TestAdvertise(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("BadOptionFails", func(t *testing.T) {
		t.Parallel()

		c := testutil.NewClock(ctrl, 0, nil)
		c.EXPECT().
			After(netsim.DefaultTTL, gomock.All()).
			Times(1)

		d := netsim.DiscoveryService{
			NS:   namespace.New(c),
			Topo: netsim.SelectAll{},
			Info: testutil.RandInfo(),
		}
		ttl, err := d.Advertise(context.Background(), "",
			func(*discovery.Options) error { return errors.New("test") })
		require.EqualError(t, err, "test")
		require.Zero(t, ttl)
	})

	t.Run("DefaultTTL", func(t *testing.T) {
		t.Parallel()

		c := testutil.NewClock(ctrl, 0, nil)
		c.EXPECT().
			After(netsim.DefaultTTL, gomock.All()).
			Times(1)

		ns := namespace.New(c)
		pi := testutil.RandInfo()
		d := netsim.DiscoveryService{
			NS:   ns,
			Topo: netsim.SelectAll{},
			Info: pi,
		}

		ttl, err := d.Advertise(ctx, "")
		require.NoError(t, err)
		require.Equal(t, netsim.DefaultTTL, ttl)

		got := ns.LoadOrCreate("").
			Peers().
			Filter(func(info *peer.AddrInfo) bool { return info.ID == pi.ID })

		require.ElementsMatch(t, netsim.InfoSlice{pi}, got)
	})

	t.Run("TTL=1000ms", func(t *testing.T) {
		t.Parallel()

		const customTTL = time.Second

		c := testutil.NewClock(ctrl, 0, nil)
		c.EXPECT().
			After(customTTL, gomock.All()).
			Times(1)

		ns := namespace.New(c)
		pi := testutil.RandInfo()
		d := netsim.DiscoveryService{
			NS:   ns,
			Topo: netsim.SelectAll{},
			Info: pi,
		}

		ttl, err := d.Advertise(ctx, "", discovery.TTL(customTTL))
		require.NoError(t, err)
		require.Equal(t, customTTL, ttl)

		got := ns.LoadOrCreate("").
			Peers().
			Filter(func(info *peer.AddrInfo) bool { return info.ID == pi.ID })

		require.ElementsMatch(t, netsim.InfoSlice{pi}, got)
	})
}

func TestFindPeers(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("DefaultOptionErrorFails", func(t *testing.T) {
		t.Parallel()

		c := testutil.NewClock(ctrl, 0, nil)
		c.EXPECT().
			After(netsim.DefaultTTL, gomock.All()).
			Times(1)

		d := netsim.DiscoveryService{
			NS:   namespace.New(c),
			Topo: failDefaultOptions{},
			Info: testutil.RandInfo(),
		}
		peers, err := d.FindPeers(context.Background(), "")
		require.EqualError(t, err, "test")
		require.Nil(t, peers)
	})

	t.Run("BadOptionFails", func(t *testing.T) {
		t.Parallel()

		c := testutil.NewClock(ctrl, 0, nil)
		c.EXPECT().
			After(netsim.DefaultTTL, gomock.All()).
			Times(1)

		d := netsim.DiscoveryService{
			NS:   namespace.New(c),
			Topo: netsim.SelectAll{},
			Info: testutil.RandInfo(),
		}
		peers, err := d.FindPeers(context.Background(), "",
			func(*discovery.Options) error { return errors.New("test") })
		require.EqualError(t, err, "test")
		require.Nil(t, peers)
	})

	t.Run("ValidationErrorFails", func(t *testing.T) {
		t.Parallel()

		c := testutil.NewClock(ctrl, 0, nil)
		c.EXPECT().
			After(netsim.DefaultTTL, gomock.All()).
			Times(1)

		d := netsim.DiscoveryService{
			NS:   namespace.New(c),
			Topo: failValidaton{},
			Info: testutil.RandInfo(),
		}
		peers, err := d.FindPeers(context.Background(), "")
		require.EqualError(t, err, "test")
		require.Nil(t, peers)
	})

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		t.Helper()

		t.Run("FoundPeers", func(t *testing.T) {
			t.Parallel()

			c := testutil.NewClock(ctrl, 0, nil)
			c.EXPECT().
				After(netsim.DefaultTTL, gomock.All()).
				Times(n)

			ns := newTestNs(c, "", n)

			d := netsim.DiscoveryService{
				NS:   ns,
				Topo: netsim.SelectAll{},
				Info: testutil.RandInfo(),
			}

			peers, err := d.FindPeers(context.Background(), "")
			require.NoError(t, err)
			require.Len(t, peers, n)
		})

		t.Run("NoPeers", func(t *testing.T) {
			t.Parallel()

			c := testutil.NewClock(ctrl, 0, nil)
			c.EXPECT().
				After(netsim.DefaultTTL, gomock.All()).
				Times(1)

			d := netsim.DiscoveryService{
				NS:   namespace.New(c),
				Topo: netsim.SelectAll{},
				Info: testutil.RandInfo(),
			}

			peers, err := d.FindPeers(context.Background(), "")
			require.NoError(t, err)
			require.Empty(t, peers)
		})
	})
}

func newTestNs(t namespace.Timer, name string, n int) netsim.NamespaceProvider {
	ns := namespace.New(t)
	for i := 0; i < n; i++ {
		ns.LoadOrCreate(name).Upsert(testutil.RandInfo(), &discovery.Options{Ttl: netsim.DefaultTTL})
	}
	return ns
}

type failValidaton struct{ netsim.SelectAll }

func (failValidaton) Validate(*discovery.Options) error {
	return errors.New("test")
}

type failDefaultOptions struct{ netsim.SelectAll }

func (failDefaultOptions) SetDefaultOptions(*discovery.Options) error {
	return errors.New("test")
}
