package matrix

import "github.com/wetware/matrix/pkg/env"

// Option for simulation.
type Option func(*Config) error

// Config is a general-purpose container for configuration.
type Config struct {
	env   Env
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

// WithEnv sets the simulation environment.
//
// If env == nil, the default global environment is used.
func WithEnv(e Env) Option {
	if e == nil {
		e = env.Global()
	}

	return func(c *Config) error {
		c.env = e
		return nil
	}
}

func withDefault(opt []Option) []Option {
	return append([]Option{
		WithEnv(nil),
	}, opt...)
}
