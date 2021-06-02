package mx

import (
	"context"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/wetware/matrix/pkg/netsim"
	"golang.org/x/sync/errgroup"
)

type (
	MapFunc    func(ctx context.Context, sim Simulation, i int, h host.Host) error
	SelectFunc func(ctx context.Context, sim Simulation, hs HostSlice) (HostSlice, error)
	FilterFunc func(int, host.Host) bool
)

// Fail is an operation that always fails with the supplied error.
func Fail(err error) OpFunc {
	return func(context.Context, Simulation, HostSlice) (HostSlice, error) {
		return nil, err
	}
}

// Map applies f to each item in the current selection.
func Map(f MapFunc) OpFunc {
	return mapper(func(hs HostSlice, hf func(MapFunc) func(int, host.Host) error) error {
		return hs.Map(hf(f))
	})
}

// Go applies f to each item in the current selection concurrently.
func Go(f MapFunc) OpFunc {
	return mapper(func(hs HostSlice, hf func(MapFunc) func(int, host.Host) error) error {
		return hs.Go(hf(f))
	})
}

// Select performs an arbitrary operation on the current selection,
// possibly returning a new selection.
//
// HostSlice -> HostSlice
func Select(f SelectFunc) OpFunc {
	return func(ctx context.Context, sim Simulation, hs HostSlice) (HostSlice, error) {
		return f(ctx, sim, hs)
	}
}

// Filter returns a new selection that contains the elements of the
// current selection for which f(element) == true.
func Filter(f FilterFunc) OpFunc {
	return Select(func(ctx context.Context, sim Simulation, hs HostSlice) (HostSlice, error) {
		return hs.Filter(f), nil
	})
}

// Announce each host in the current selection using the supplied topology.
func Announce(t netsim.Topology, ns string, opt ...discovery.Option) OpFunc {
	return Go(func(ctx context.Context, sim Simulation, i int, h host.Host) error {
		var d = sim.NewDiscovery(h, t)
		_, err := d.Advertise(ctx, ns, opt...)
		return err
	})
}

// Discover peers for each host in the current selection using the supplied topology.
func Discover(t netsim.Topology, ns string, opt ...discovery.Option) OpFunc {
	return Go(func(ctx context.Context, sim Simulation, i int, h host.Host) (err error) {
		var (
			d  = sim.NewDiscovery(h, t)
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
	return func(ctx context.Context, sim Simulation, hs HostSlice) (HostSlice, error) {
		return hs, f(hs, func(mf MapFunc) func(int, host.Host) error {
			return func(i int, h host.Host) error {
				return mf(ctx, sim, i, h)
			}
		})
	}
}
