package netsim

import (
	"context"
	"math/rand"
	"testing"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/pkg/clock"
	"golang.org/x/sync/errgroup"
)

func TestTopology(t *testing.T) {
	t.Parallel()
	t.Helper()

	var (
		p     = newTestNs(clock.New(), "", n)
		ns    = p.LoadOrCreate("")
		local = randinfo()
	)

	t.Run("SelectAll", func(t *testing.T) {
		t.Parallel()
		t.Helper()

		t.Run("Limit", func(t *testing.T) {
			t.Parallel()

			const limit = 5

			var s SelectAll

			as, err := run(ns, s, local, discovery.Limit(limit))
			require.NoError(t, err)

			require.Subset(t, ns.Peers(), as)
			require.Len(t, as, limit)

			peers := defaultLoader{}.load(ns, local)
			require.NotContains(t, peers, local)
			require.Equal(t, peers[:limit], as)
		})
	})

	t.Run("SelectRandom", func(t *testing.T) {
		t.Parallel()
		t.Helper()

		t.Run("GlobalSource", func(t *testing.T) {
			t.Parallel()

			var s SelectRandom
			ps, err := run(ns, &s, local)
			require.NoError(t, err)

			peers := defaultLoader{}.load(ns, local)
			require.NotContains(t, peers, local)
			require.ElementsMatch(t, peers, ps)
		})

		t.Run("Reproducible", func(t *testing.T) {
			t.Parallel()

			as0, err := run(ns, &SelectRandom{
				Src: rand.NewSource(42),
			}, local)
			require.NoError(t, err)

			as1, err := run(ns, &SelectRandom{
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
			peers     = loadAllPeers(ns)
			neighbors = make(InfoSlice, len(peers))
		)

		/*
		 * 'neighbors' should be identical to 'peers', except that it is
		 *  rotated by 1 towards the tail.
		 */

		var g errgroup.Group
		for i, info := range peers {
			g.Go(func(i int, info *peer.AddrInfo) func() error {
				return func() (err error) {
					var as InfoSlice

					if as, err = run(ns, &SelectRing{}, info); err == nil {
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

func run(ns Namespace, topo Topology, info *peer.AddrInfo, opt ...discovery.Option) (InfoSlice, error) {
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

	return topo.Select(context.Background(), ns, info, opts)
}

func newOption() *discovery.Options {
	return &discovery.Options{Other: make(map[interface{}]interface{})}
}

func loadAllPeers(ns Namespace) InfoSlice {
	return defaultLoader{}.load(ns, new(peer.AddrInfo))
}
