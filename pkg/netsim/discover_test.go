package netsim

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/pkg/clock"
)

const n = 10

var t0 = time.Date(2021, 04, 9, 8, 0, 0, 0, time.UTC)

func TestAdvertise(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("BadOptionFails", func(t *testing.T) {
		t.Parallel()

		d := DiscoveryService{
			NS:   nsmap(clock.New()),
			Topo: SelectAll{},
			Info: randinfo(),
		}
		ttl, err := d.Advertise(context.Background(), "",
			func(*discovery.Options) error { return errors.New("test") })
		require.EqualError(t, err, "test")
		require.Zero(t, ttl)
	})

	t.Run("DefaultTTL", func(t *testing.T) {
		t.Parallel()

		c := clock.New()
		c.Advance(t0)

		ns := nsmap(c)
		p := randinfo()
		d := DiscoveryService{
			NS:   ns,
			Topo: SelectAll{},
			Info: p,
		}

		t.Run("Advertise", func(t *testing.T) {
			ttl, err := d.Advertise(ctx, "")
			require.NoError(t, err)
			require.Equal(t, DefaultTTL, ttl)

			got := ns.LoadOrCreate("").
				Peers().
				Filter(func(info *peer.AddrInfo) bool { return info.ID == p.ID })

			require.ElementsMatch(t, InfoSlice{p}, got)
		})

		t.Run("Expire", func(t *testing.T) {
			c.Advance(t0.Add(DefaultTTL + c.Accuracy()))

			require.Eventually(t, func() bool {
				return len(ns.LoadOrCreate("").Peers()) == 0
			}, time.Millisecond*100, time.Millisecond*10,
				"peer was not expired after %s", DefaultTTL)
		})
	})

	t.Run("TTL=1000ms", func(t *testing.T) {
		t.Parallel()

		const customTTL = time.Second
		c := clock.New()
		c.Advance(t0)

		ns := nsmap(c)
		p := randinfo()
		d := DiscoveryService{
			NS:   ns,
			Topo: SelectAll{},
			Info: p,
		}

		t.Run("Advertise", func(t *testing.T) {
			ttl, err := d.Advertise(ctx, "", discovery.TTL(customTTL))
			require.NoError(t, err)
			require.Equal(t, customTTL, ttl)

			got := ns.LoadOrCreate("").
				Peers().
				Filter(func(info *peer.AddrInfo) bool { return info.ID == p.ID })

			require.ElementsMatch(t, InfoSlice{p}, got)
		})

		t.Run("Expire", func(t *testing.T) {
			c.Advance(t0.Add(customTTL + c.Accuracy()))

			require.Eventually(t, func() bool {
				return len(ns.LoadOrCreate("").Peers()) == 0
			}, time.Millisecond*100, time.Millisecond*10,
				"peer was not expired after %s", customTTL)
		})
	})
}

func TestFindPeers(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("DefaultOptionErrorFails", func(t *testing.T) {
		t.Parallel()

		d := DiscoveryService{
			NS:   nsmap(clock.New()),
			Topo: failDefaultOptions{},
			Info: randinfo(),
		}
		peers, err := d.FindPeers(context.Background(), "")
		require.EqualError(t, err, "test")
		require.Nil(t, peers)
	})

	t.Run("BadOptionFails", func(t *testing.T) {
		t.Parallel()

		d := DiscoveryService{
			NS:   nsmap(clock.New()),
			Topo: SelectAll{},
			Info: randinfo(),
		}
		peers, err := d.FindPeers(context.Background(), "",
			func(*discovery.Options) error { return errors.New("test") })
		require.EqualError(t, err, "test")
		require.Nil(t, peers)
	})

	t.Run("ValidationErrorFails", func(t *testing.T) {
		t.Parallel()

		d := DiscoveryService{
			NS:   nsmap(clock.New()),
			Topo: failValidaton{},
			Info: randinfo(),
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

			ns := newTestNs(clock.New(), "", n)

			d := DiscoveryService{
				NS:   ns,
				Topo: SelectAll{},
				Info: randinfo(),
			}

			peers, err := d.FindPeers(context.Background(), "")
			require.NoError(t, err)
			require.Len(t, peers, n)
		})

		t.Run("NoPeers", func(t *testing.T) {
			t.Parallel()

			d := DiscoveryService{
				NS:   nsmap(clock.New()),
				Topo: SelectAll{},
				Info: randinfo(),
			}

			peers, err := d.FindPeers(context.Background(), "")
			require.NoError(t, err)
			require.Empty(t, peers)
		})
	})
}

func newTestNs(t Timer, ns string, n int) NamespaceProvider {
	p := nsmap(t)
	for i := 0; i < n; i++ {
		p.LoadOrCreate(ns).Upsert(randinfo(), &discovery.Options{Ttl: DefaultTTL})
	}
	return p
}

func randinfo() *peer.AddrInfo {
	id := randID()
	return &peer.AddrInfo{
		ID:    id,
		Addrs: []multiaddr.Multiaddr{newAddr(id)},
	}
}

func newAddr(id peer.ID) multiaddr.Multiaddr {
	ma, err := inproc.ResolveString("/inproc/~")
	if err != nil {
		panic(err)
	}

	return ma.Encapsulate(multiaddr.StringCast(fmt.Sprintf("/p2p/%s", id)))
}

func randID() peer.ID {
	return newID(randStr(5))
}

func randStr(n int) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	b := make([]rune, n)
	for i := range b {
		b[i] = rune(alphabet[rand.Intn(len(alphabet))])
	}

	return string(b)
}

func hash(b []byte) []byte {
	h, _ := multihash.Sum(b, multihash.SHA2_256, -1)
	return []byte(h)
}

func newID(s string) peer.ID {
	id, err := peer.Decode(base58.Encode(hash([]byte(s))))
	if err != nil {
		panic(err)
	}

	return id
}

type failValidaton struct{ SelectAll }

func (failValidaton) Validate(*discovery.Options) error {
	return errors.New("test")
}

type failDefaultOptions struct{ SelectAll }

func (failDefaultOptions) SetDefaultOptions(*discovery.Options) error {
	return errors.New("test")
}
