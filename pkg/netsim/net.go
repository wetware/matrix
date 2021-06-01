package netsim

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/config"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/wetware/matrix/pkg/clock"
)

type Timer interface {
	After(d time.Duration, callback func()) clock.CancelFunc
}

type Env struct {
	inproc inproc.Env
	NS     NamespaceProvider
}

func New(t Timer) *Env {
	return &Env{
		inproc: inproc.NewEnv(),
		NS:     nsmap(t),
	}
}

// NewHost assembles and creates a new libp2p host that uses the
// simulation's network.
//
// Env configures hosts to use an in-process network and therefore
// overrides the following options:
//
// - libp2p.Transport
// - libp2p.NoTransports
// - libp2p.ListenAddr
// - libp2p.ListenAddrStrings
// - libp2p.NoListenAddrs
//
// Users SHOULD NOT pass any of the above options to NewHost.
func (env Env) NewHost(ctx context.Context, opt []config.Option) (host.Host, error) {
	return libp2p.New(ctx, env.hostopt(opt)...)
}

func (env Env) hostopt(opt []libp2p.Option) []config.Option {
	return append(opt,
		// override options that users may have passed
		libp2p.NoListenAddrs,
		libp2p.ListenAddrStrings("/inproc/~"),
		libp2p.NoTransports,
		libp2p.Transport(inproc.New(inproc.WithEnv(env.inproc))))
}
