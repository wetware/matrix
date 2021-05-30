package env

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/discover"
)

// Env encapsulates bindings in an isolated address space.
type Env struct {
	clock *clock.Clock
	net   inproc.Env
	ns    *nsMap
}

func New(ctx context.Context, opt ...Option) *Env {
	env := &Env{net: inproc.NewEnv()}
	defer run(ctx, env)

	for _, option := range withDefault(opt) {
		option(env)
	}

	env.ns = nsmap(env.clock)

	return env
}

func (env *Env) Network() inproc.Env { return env.net }
func (env *Env) Clock() *Clock       { return (*Clock)(env.clock) }

func (env *Env) NewDiscovery(info peer.AddrInfo, s discover.Strategy) discovery.Discovery {
	return discover.New(env.ns, s, &info)
}

func run(ctx context.Context, env *Env) {
	go func() {
		ticker := time.NewTicker(env.clock.Accuracy())
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				env.clock.Advance(t)

			case <-ctx.Done():
				return
			}
		}
	}()
}

type Clock clock.Clock

func (c *Clock) Accuracy() time.Duration { return (*clock.Clock)(c).Accuracy() }

func (c *Clock) After(d time.Duration, callback func()) clock.CancelFunc {
	return (*clock.Clock)(c).After(d, callback)
}

func (c *Clock) Ticker(d time.Duration, callback func()) clock.CancelFunc {
	return (*clock.Clock)(c).Ticker(d, callback)
}
