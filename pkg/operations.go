package matrix

import (
	"context"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/wetware/matrix/pkg/discover"
	"golang.org/x/sync/errgroup"
)

type Op struct {
	env Env
	f   OpFunc
}

func (op Op) Call(ctx context.Context, hs ...host.Host) error {
	return op.f(ctx, op.env, hs)
}

func (op Op) Must(ctx context.Context, hs ...host.Host) {
	if err := op.Call(ctx, hs...); err != nil {
		panic(err)
	}
}

func (op Op) Then(f OpFunc) Op {
	return Op{
		env: op.env,
		f: func(ctx context.Context, env Env, hs HostSlice) error {
			if err := op.Call(ctx, hs...); err != nil {
				return err
			}

			return f(ctx, env, hs)
		},
	}
}

type OpFunc func(ctx context.Context, env Env, hs HostSlice) error

func Announce(s discover.Strategy, ns string, opt ...discovery.Option) OpFunc {
	return func(ctx context.Context, env Env, hs HostSlice) error {
		// NewDiscovery(info peer.AddrInfo, s discover.Strategy) discovery.Discovery
		return hs.Go(func(_ int, h host.Host) error {
			_, err := env.NewDiscovery(*host.InfoFromHost(h), s).
				Advertise(ctx, ns, opt...)
			return err
		})
	}
}

func Discover(s discover.Strategy, ns string, opt ...discovery.Option) OpFunc {
	return func(ctx context.Context, env Env, hs HostSlice) error {
		return hs.Go(func(_ int, h host.Host) error {
			ps, err := env.NewDiscovery(*host.InfoFromHost(h), s).
				FindPeers(ctx, ns, opt...)
			if err != nil {
				return nil
			}

			var g errgroup.Group
			for info := range ps {
				g.Go(func(info peer.AddrInfo) func() error {
					return func() error {
						return h.Connect(ctx, info)
					}
				}(info))
			}
			return g.Wait()
		})
	}
}

type HostSlice []host.Host

func (hs HostSlice) Len() int           { return len(hs) }
func (hs HostSlice) Less(i, j int) bool { return hs[i].ID() < hs[j].ID() }
func (hs HostSlice) Swap(i, j int)      { hs[i], hs[j] = hs[j], hs[i] }

func (hs HostSlice) Apply(f HostFunc) (err error) {
	for i, h := range hs {
		if err = f(i, h); err != nil {
			break
		}
	}

	return
}

func (hs HostSlice) Go(f HostFunc) error {
	var g errgroup.Group
	for i, h := range hs {
		g.Go(func(i int, h host.Host) func() error {
			return func() (err error) {
				err = f(i, h)
				return
			}
		}(i, h))
	}
	return g.Wait()
}

// HostFunc can modify a host.
type HostFunc func(int, host.Host) error

// func Topology(s discover.Strategy, ns string, opt ...discovery.Option) Op {
// 	return Discover(s, ns, opt...).
// 		After(Announce(s, ns, opt...))
// }

// func Announce(s discover.Strategy, ns string, opt ...discovery.Option) Op {
// 	return func(ctx context.Context, env *env.Env, hs HostSlice) (HostSlice, error) {
// 		return hs.Go(func(_ int, h host.Host) error {
// 			_, err := env.NewDiscovery(*host.InfoFromHost(h), s).
// 				Advertise(ctx, ns, opt...)
// 			return err
// 		})
// 	}
// }

// func Discover(s discover.Strategy, ns string, opt ...discovery.Option) Op {
// 	return func(ctx context.Context, env *env.Env, h HostSlice) (HostSlice, error) {
// 		return h.Go(func(_ int, h host.Host) error {
// 			ps, err := env.NewDiscovery(*host.InfoFromHost(h), s).
// 				FindPeers(ctx, ns, opt...)
// 			if err != nil {
// 				return nil
// 			}

// 			var g errgroup.Group
// 			for info := range ps {
// 				g.Go(func(info peer.AddrInfo) func() error {
// 					return func() error {
// 						return h.Connect(ctx, info)
// 					}
// 				}(info))
// 			}
// 			return g.Wait()
// 		})
// 	}
// }
