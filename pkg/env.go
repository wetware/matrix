package matrix

import (
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/config"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/wetware/matrix/pkg/discover"
	"github.com/wetware/matrix/pkg/env"
)

type Env interface {
	Clock() *env.Clock
	Network() inproc.Env
	NewHost(ctx context.Context, opt ...Option) (host.Host, error)
	MustHost(ctx context.Context, opt ...Option) host.Host
	NewDiscovery(info peer.AddrInfo, s discover.Strategy) discovery.Discovery
	Op(...OpFunc) Op
}

type environment struct{ *env.Env }

func New(ctx context.Context) Env { return environment{env.New(ctx)} }

func (env environment) Op(ops ...OpFunc) Op {
	var of OpFunc
	for _, op := range ops {
		of = of.Then(op)
	}

	return Op{env: env, call: of}
}

// NewHost returns a new host with a random identity and the default
// in-process transport.
func (env environment) NewHost(ctx context.Context, opt ...Option) (host.Host, error) {
	cfg, err := options(opt)
	if err != nil {
		return nil, err
	}

	return cfg.newHost(ctx, env)
}

// MustHost returns a host or panics if an error was encountered.
func (env environment) MustHost(ctx context.Context, opt ...Option) host.Host {
	h, err := env.NewHost(ctx, opt...)
	must(err)
	return h
}

func transport(ctx context.Context, e inproc.Env) config.Option {
	return libp2p.Transport(inproc.New(inproc.WithEnv(e)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
