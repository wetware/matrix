package mx

import (
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/namespace"
	"github.com/wetware/matrix/pkg/netsim"
)

type Option func(sim *Simulation)

// WithClock sets the simulation clock.  If c == nil, a default
// clock with 10ms accuracy is used.
func WithClock(c ClockController) Option {
	if c == nil {
		c = clock.New()
	}

	return func(sim *Simulation) {
		sim.c = c
	}
}

// WithHostFactory sets the host factory for the simulation.
// If f == nil, a default factory is used that constructs
// hosts using an in-process transport.
func WithHostFactory(f HostFactory) Option {
	if f == nil {
		f = &netsim.HostFactory{}
	}

	return func(sim *Simulation) {
		sim.h = f
	}
}

// WithNamespaceFactory sets the namespace implementation
// for the simulation.  If ns == nil, a default namespace
// implementation is used.
func WithNamespaceFactory(f func(Clock) netsim.NamespaceProvider) Option {
	if f == nil {
		f = func(c Clock) netsim.NamespaceProvider { return namespace.New(c) }
	}

	return func(sim *Simulation) {
		sim.init = f
	}
}

func withDefault(opt []Option) []Option {
	return append([]Option{
		WithClock(nil),
		WithHostFactory(nil),
		WithNamespaceFactory(nil),
	}, finalize(opt)...)
}

func finalize(opt []Option) []Option {
	return append(opt, func(sim *Simulation) {
		sim.ns = sim.init(sim.c)
	})
}
