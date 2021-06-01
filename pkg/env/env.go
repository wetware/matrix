package env

import (
	"context"
	"time"

	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/wetware/matrix/pkg/clock"
)

// Env encapsulates bindings in an isolated address space.
type Env struct {
	clock *clock.Clock
	net   inproc.Env
}

func New(ctx context.Context, opt ...Option) *Env {
	env := &Env{net: inproc.NewEnv()}
	defer run(ctx, env)

	for _, option := range withDefault(opt) {
		option(env)
	}

	return env
}

func (env *Env) Network() inproc.Env { return env.net }

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
