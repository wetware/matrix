package matrix

import (
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/config"
)

// Option for simulation.
type Option func(*Config) error

// Config is a general-purpose container for configuration.
type Config struct {
	other map[interface{}]interface{}
}

func options(opt []Option) (*Config, error) {
	var c Config
	for _, option := range withDefault(opt) {
		if err := option(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c *Config) SetCustom(k, v interface{}) {
	if c.other == nil {
		c.other = make(map[interface{}]interface{})
	}

	c.other[k] = v
}

func (c *Config) GetCustom(k interface{}) (v interface{}, ok bool) {
	if c.other != nil {
		v, ok = c.other[k]
	}

	return
}

func (c *Config) newHost(ctx context.Context, env Env) (host.Host, error) {
	opt, _ := c.other[keyHostOpt].([]libp2p.Option)
	return libp2p.New(ctx, c.hostopt(ctx, env, opt)...)
}

func (c *Config) hostopt(ctx context.Context, env Env, opt []libp2p.Option) []config.Option {
	return append([]config.Option{
		transport(ctx, env.Network()),
		libp2p.NoListenAddrs,
		libp2p.ListenAddrStrings("/inproc/~"),
	}, opt...)
}

func withDefault(opt []Option) []Option {
	return append([]Option{
		// ...
	}, opt...)
}

type key uint8

const (
	keyHostOpt key = iota
)
