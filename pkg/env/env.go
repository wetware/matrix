package env

import (
	"context"
	"sync"
	"time"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/discover"
)

var (
	once      sync.Once
	globalEnv *Env
)

func Global() *Env {
	once.Do(func() { globalEnv = New(context.Background()) })
	return globalEnv
}

// Env encapsulates bindings in an isolated address space.
type Env struct {
	clock *clock.Clock
	proc  goprocess.Process
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

func (env *Env) Clock() *Clock              { return (*Clock)(env.clock) }
func (env *Env) Process() goprocess.Process { return env.proc }
func (env *Env) Network() inproc.Env        { return env.net }

func (env *Env) NewDiscovery(info peer.AddrInfo, s discover.Strategy) discovery.Discovery {
	return discover.New(env.ns, s, &info)
}

func run(ctx context.Context, env *Env) {
	env.proc = process(ctx).Go(func(p goprocess.Process) {
		ticker := time.NewTicker(env.clock.Accuracy())
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				env.clock.Advance(t)

			case <-p.Closing():
				return
			}
		}
	})
}

func process(ctx context.Context) goprocess.Process {
	if s, ok := ctx.(interface{ String() string }); ok {
		switch s.String() {
		case "context.TODO", "context.Background":
			return goprocess.Background()
		}
	}

	return goprocessctx.WithContext(ctx)
}

type Clock clock.Clock

func (c *Clock) Accuracy() time.Duration { return (*clock.Clock)(c).Accuracy() }

func (c *Clock) After(d time.Duration, callback func()) clock.CancelFunc {
	return (*clock.Clock)(c).After(d, callback)
}

func (c *Clock) Ticker(d time.Duration, callback func()) clock.CancelFunc {
	return (*clock.Clock)(c).Ticker(d, callback)
}