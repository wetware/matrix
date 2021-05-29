package env

import "github.com/wetware/matrix/pkg/clock"

type Option func(env *Env)

func WithClock(c *clock.Clock) Option {
	if c == nil {
		c = clock.New()
	}

	return func(env *Env) {
		env.clock = c
	}
}

func withDefault(opt []Option) []Option {
	return append([]Option{
		WithClock(nil),
	}, opt...)
}
