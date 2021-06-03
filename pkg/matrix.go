//go:generate mockgen -destination ../internal/mock/pkg/matrix/matrix.go github.com/wetware/matrix/pkg ClockController,HostFactory

package mx

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/config"
	"github.com/wetware/matrix/pkg/netsim"
)

type Clock interface {
	Accuracy() time.Duration
	After(d time.Duration, callback func()) (cancel func())
	Ticker(userExpire time.Duration, callback func()) (cancel func())
}

type ClockController interface {
	Clock
	Advance(time.Time)
}

type HostFactory interface {
	NewHost(ctx context.Context, opt []config.Option) (host.Host, error)
}

type Simulation struct {
	c ClockController
	h HostFactory

	init func(c Clock) netsim.NamespaceProvider
	ns   netsim.NamespaceProvider
}

func New(ctx context.Context, opt ...Option) Simulation {
	var sim Simulation
	for _, option := range withDefault(opt) {
		option(&sim)
	}

	go sim.run(ctx)

	return sim
}

func (s Simulation) Clock() Clock { return s.c }

// NewHost assembles and creates a new libp2p host that uses the
// simulation's network.
//
// The simulation configures hosts to use an in-process network,
// overriding the following options:
//
// - libp2p.Transport
// - libp2p.NoTransports
// - libp2p.ListenAddr
// - libp2p.ListenAddrStrings
// - libp2p.NoListenAddrs
//
// Users SHOULD NOT pass any of the above options to NewHost.
func (s Simulation) NewHost(ctx context.Context, opt ...config.Option) (host.Host, error) {
	return s.h.NewHost(ctx, opt)
}

// MustHost returns a host or panics if an error was encountered.
func (s Simulation) MustHost(ctx context.Context, opt ...config.Option) host.Host {
	h, err := s.NewHost(ctx, opt...)
	must(err)
	return h
}

// NewHostSet builds and configures n hosts with identical parameters.
//
// See NewHost.
func (s Simulation) NewHostSet(ctx context.Context, n int, opt ...config.Option) (Selection, error) {
	hs := make(Selection, n)
	return hs, hs.Go(func(i int, _ host.Host) (err error) {
		hs[i], err = s.h.NewHost(ctx, opt)
		return
	})
}

// MustHostSet calls NewHostSet with the supplied parameters and panics if
// an error is encountered.
func (s Simulation) MustHostSet(ctx context.Context, n int, opt ...config.Option) Selection {
	hs, err := s.NewHostSet(ctx, n, opt...)
	must(err)
	return hs
}

// NewDiscovery returns a discovery.Discovery implementation that
// supports the Simulation's in-process network.
//
// The topology parameter t can be used to specify an initial
// connection topology.  All peers must use the same instance
// of t in order to obtain the desired topology.
//
// If t == nil, the topology defaults to netsim.SelectAll.
func (s Simulation) NewDiscovery(h host.Host, t netsim.Topology) *netsim.DiscoveryService {
	return &netsim.DiscoveryService{
		NS:   s.ns,
		Info: host.InfoFromHost(h),
		Topo: topology(t),
	}
}

func (s Simulation) run(ctx context.Context) {
	ticker := time.NewTicker(s.c.Accuracy())
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			s.c.Advance(t)

		case <-ctx.Done():
			return
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func topology(t netsim.Topology) netsim.Topology {
	if t != nil {
		return t
	}

	return netsim.SelectAll{}
}
