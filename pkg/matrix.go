package matrix

import (
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/config"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
)

type Env interface {
	Network() inproc.Env
}

// NewHost
func NewHost(ctx context.Context, opt ...Option) (host.Host, error) {
	ho, err := hostopt(opt)
	if err != nil {
		return nil, err
	}

	return libp2p.New(ctx, ho...)
}

func hostopt(opt []Option) ([]config.Option, error) {
	cfg, err := options(opt)
	if err != nil {
		return nil, err
	}

	return []config.Option{
		transport(cfg),
		libp2p.NoListenAddrs,
		libp2p.ListenAddrStrings("/inproc/~"),
	}, nil

}

func transport(cfg *Config) config.Option {
	return libp2p.Transport(inproc.New(inproc.WithEnv(cfg.env.Network())))
}
