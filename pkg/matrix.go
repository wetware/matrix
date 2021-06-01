package mx

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/config"
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/net"
)

type Clock interface {
	Accuracy() time.Duration
	After(d time.Duration, callback func()) clock.CancelFunc
	Ticker(userExpire time.Duration, callback func()) clock.CancelFunc
}

type Simulation interface {
	Clock() Clock
	NewHost(ctx context.Context, opt ...config.Option) (host.Host, error)
	MustHost(ctx context.Context, opt ...config.Option) host.Host
	NewDiscovery(info peer.AddrInfo, s net.Topology) discovery.Discovery
	Op(...OpFunc) Op
}

type sim struct {
	*net.Env
	c *clock.Clock
}

func New(ctx context.Context) Simulation {
	c := clock.New()
	go tick(ctx, c)

	return sim{
		Env: net.New(c),
		c:   c,
	}
}

func (s sim) Clock() Clock { return s.c }

func (s sim) Op(ops ...OpFunc) Op {
	var of OpFunc
	for _, op := range ops {
		of = of.Then(op)
	}

	return Op{sim: s, call: of}
}

// MustHost returns a host or panics if an error was encountered.
func (s sim) MustHost(ctx context.Context, opt ...config.Option) host.Host {
	h, err := s.NewHost(ctx, opt...)
	must(err)
	return h
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func tick(ctx context.Context, c *clock.Clock) {
	ticker := time.NewTicker(c.Accuracy())
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			c.Advance(t)

		case <-ctx.Done():
			return
		}
	}
}
