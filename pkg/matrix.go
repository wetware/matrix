package matrix

import (
	"context"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/wetware/matrix/pkg/discover"
	"golang.org/x/sync/errgroup"
)

type (
	HostSlice []host.Host

	MapFunc    func(ctx context.Context, env Env, i int, h host.Host) error
	SelectFunc func(ctx context.Context, env Env, hs HostSlice) (HostSlice, error)

	FilterFunc func(int, host.Host) bool
)

func Map(f MapFunc) OpFunc {
	return mapper(func(hs HostSlice, hf func(MapFunc) func(int, host.Host) error) error {
		return hs.Map(hf(f))
	})
}

func Go(f MapFunc) OpFunc {
	return mapper(func(hs HostSlice, hf func(MapFunc) func(int, host.Host) error) error {
		return hs.Go(hf(f))
	})
}

func Select(f SelectFunc) OpFunc {
	return func(env Env) func(ctx context.Context) Maybe {
		return func(ctx context.Context) Maybe {
			return func(hs HostSlice) (HostSlice, error) {
				return f(ctx, env, hs)
			}
		}
	}
}

func Filter(f FilterFunc) OpFunc {
	return Select(func(ctx context.Context, env Env, hs HostSlice) (HostSlice, error) {
		sel := hs[:0]
		for i, h := range hs {
			if f(i, h) {
				sel = append(sel, h)
			}
		}
		return sel, nil
	})
}

func Announce(s discover.Strategy, ns string, opt ...discovery.Option) OpFunc {
	return Go(func(ctx context.Context, env Env, i int, h host.Host) error {
		var d = env.NewDiscovery(*host.InfoFromHost(h), s)
		_, err := d.Advertise(ctx, ns, opt...)
		return err
	})
}

func Discover(s discover.Strategy, ns string, opt ...discovery.Option) OpFunc {
	return Go(func(ctx context.Context, env Env, i int, h host.Host) (err error) {
		var (
			d  = env.NewDiscovery(*host.InfoFromHost(h), s)
			g  errgroup.Group
			ps <-chan peer.AddrInfo
		)

		if ps, err = d.FindPeers(ctx, ns, opt...); err != nil {
			return err
		}

		for info := range ps {
			if info.ID != h.ID() {
				g.Go(connect(ctx, h, info))
			}
		}
		return g.Wait()
	})
}

func connect(ctx context.Context, h host.Host, info peer.AddrInfo) func() error {
	return func() error {
		return h.Connect(ctx, info)
	}
}

func mapper(f func(hs HostSlice, hf func(MapFunc) func(int, host.Host) error) error) OpFunc {
	return func(env Env) func(ctx context.Context) Maybe {
		return func(ctx context.Context) Maybe {
			return func(hs HostSlice) (HostSlice, error) {
				return hs, f(hs, func(mf MapFunc) func(int, host.Host) error {
					return func(i int, h host.Host) error {
						return mf(ctx, env, i, h)
					}
				})
			}
		}
	}
}
