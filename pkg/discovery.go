package mx

import (
	"context"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/wetware/matrix/pkg/netsim"
	"golang.org/x/sync/errgroup"
)

// Announce each host in the current selection using the supplied topology.
func Announce(sim Simulation, t netsim.Topology, ns string, opt ...discovery.Option) MapFunc {
	return func(ctx context.Context, _ int, h host.Host) (err error) {
		_, err = sim.NewDiscovery(h, t).Advertise(ctx, ns, opt...)
		return
	}
}

// Discover peers for each host in the current selection using the supplied topology.
func Discover(sim Simulation, t netsim.Topology, ns string, opt ...discovery.Option) MapFunc {
	return func(ctx context.Context, _ int, h host.Host) error {
		ps, err := sim.NewDiscovery(h, t).FindPeers(ctx, ns, opt...)
		if err != nil {
			return err
		}

		var g errgroup.Group
		for info := range ps {
			if info.ID != h.ID() {
				g.Go(connect(ctx, h, info))
			}
		}

		return g.Wait()
	}
}

func Topology(sim Simulation, t netsim.Topology, ns string, opt ...discovery.Option) Op {
	return Go(func(ctx context.Context, i int, h host.Host) error {
		build := Announce(sim, t, ns, opt...).
			Then(Discover(sim, t, ns, opt...))
		return build(ctx, i, h)
	})
}
