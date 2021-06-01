package netsim_test

import (
	"context"
	"math/rand"
	"sort"
	"testing"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/internal/testutil"
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/netsim"
	"golang.org/x/sync/errgroup"
)

func TestTopology(t *testing.T) {
	t.Parallel()
	t.Helper()

	var (
		p     = newTestNs(clock.New(), "", n)
		s     = p.LoadOrCreate("")
		local = testutil.RandInfo()
	)

	t.Run("SelectAll", func(t *testing.T) {
		t.Parallel()
		t.Helper()

		t.Run("Limit", func(t *testing.T) {
			t.Parallel()

			const limit = 5

			var topo netsim.SelectAll

			as, err := run(s, topo, local, discovery.Limit(limit))
			require.NoError(t, err)

			require.Subset(t, s.Peers(), as)
			require.Len(t, as, limit)

			peers := load(s, local)
			require.NotContains(t, peers, local)
			require.Equal(t, peers[:limit], as)
		})
	})

	t.Run("SelectRandom", func(t *testing.T) {
		t.Parallel()
		t.Helper()

		t.Run("GlobalSource", func(t *testing.T) {
			t.Parallel()

			var topo netsim.SelectRandom
			ps, err := run(s, &topo, local)
			require.NoError(t, err)

			peers := load(s, local)
			require.NotContains(t, peers, local)
			require.ElementsMatch(t, peers, ps)
		})

		t.Run("Reproducible", func(t *testing.T) {
			t.Parallel()

			as0, err := run(s, &netsim.SelectRandom{
				Src: rand.NewSource(42),
			}, local)
			require.NoError(t, err)

			as1, err := run(s, &netsim.SelectRandom{
				Src: rand.NewSource(42),
			}, local)
			require.NoError(t, err)

			require.ElementsMatch(t, as0, as1)
			require.Equal(t, as0, as1)
		})
	})

	t.Run("SelectRing", func(t *testing.T) {
		t.Parallel()

		var (
			peers     = loadAllPeers(s)
			neighbors = make(netsim.InfoSlice, len(peers))
		)

		/*
		 * 'neighbors' should be identical to 'peers', except that it is
		 *  rotated by 1 towards the tail.
		 */

		var g errgroup.Group
		for i, info := range peers {
			g.Go(func(i int, info *peer.AddrInfo) func() error {
				return func() (err error) {
					var as netsim.InfoSlice

					if as, err = run(s, &netsim.SelectRing{}, info); err == nil {
						neighbors[i] = as[0]
					}

					return
				}
			}(i, info))
		}

		require.NoError(t, g.Wait())
		require.Equal(t, peers,
			// same array, except that tail is appended to head.
			append(neighbors[n-1:], neighbors[:n-1]...))
	})
}

func run(s netsim.Scope, topo netsim.Topology, info *peer.AddrInfo, opt ...discovery.Option) (netsim.InfoSlice, error) {
	opts := newOption()
	if err := topo.SetDefaultOptions(opts); err != nil {
		return nil, err
	}

	for _, option := range opt {
		if err := option(opts); err != nil {
			return nil, err
		}
	}

	// validate?
	if v, ok := topo.(interface {
		Validate(*discovery.Options) error
	}); ok {
		if err := v.Validate(opts); err != nil {
			return nil, err
		}
	}

	return topo.Select(context.Background(), s, info, opts)
}

func newOption() *discovery.Options {
	return &discovery.Options{Other: make(map[interface{}]interface{})}
}

func loadAllPeers(s netsim.Scope) netsim.InfoSlice {
	return load(s, new(peer.AddrInfo))
}

func load(ps interface{ Peers() netsim.InfoSlice }, local *peer.AddrInfo) netsim.InfoSlice {
	is := ps.Peers()
	sort.Sort(is)
	return is.
		Filter(func(info *peer.AddrInfo) bool { return info.ID != local.ID })
}
