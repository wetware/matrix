package net

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/config"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/wetware/matrix/pkg/clock"
)

type Timer interface {
	After(d time.Duration, callback func()) clock.CancelFunc
}

type Env struct {
	inproc inproc.Env
	ns     *nsMap
}

func New(t Timer) *Env {
	return &Env{
		inproc: inproc.NewEnv(),
		ns:     nsmap(t),
	}
}

func (env Env) NewDiscovery(info peer.AddrInfo, t Topology) discovery.Discovery {
	d := discoveryService{
		ns:       env.ns,
		info:     &info,
		topo:     t,
		validate: func(*discovery.Options) error { return nil },
	}

	if t == nil {
		t = SelectAll{}
	}

	if v, ok := t.(validator); ok {
		d.validate = v.Validate
	}

	return d
}

// NewHost returns a new host with a random identity and the default
// in-process transport.
func (env Env) NewHost(ctx context.Context, opt ...config.Option) (host.Host, error) {
	return libp2p.New(ctx, env.hostopt(opt)...)
}

func (env Env) hostopt(opt []libp2p.Option) []config.Option {
	return append([]config.Option{
		libp2p.NoTransports,
		libp2p.Transport(inproc.New(inproc.WithEnv(env.inproc))),
		libp2p.NoListenAddrs,
		libp2p.ListenAddrStrings("/inproc/~"),
	}, opt...)
}
