package mx

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"golang.org/x/sync/errgroup"
)

type Maybe func(HostSlice) (HostSlice, error)

func (m Maybe) Then(next Maybe) Maybe {
	if m == nil {
		return next
	}

	return func(hs HostSlice) (_ HostSlice, err error) {
		if hs, err = m(hs); err == nil {
			return next(hs)
		}

		return
	}
}

type OpFunc func(s Simulation) func(ctx context.Context) Maybe

func (fn OpFunc) Then(next OpFunc) OpFunc {
	if fn == nil {
		return next
	}

	return func(s Simulation) func(context.Context) Maybe {
		return func(ctx context.Context) Maybe {
			return fn(s)(ctx).Then(next(s)(ctx))
		}
	}
}

type Op struct {
	sim  Simulation
	call OpFunc
}

func (op Op) Then(call OpFunc) Op {
	return Op{sim: op.sim, call: op.call.Then(call)}
}

func (op Op) Call(ctx context.Context, hs ...host.Host) OpCall {
	return func() error {
		_, err := op.call(op.sim)(ctx)(hs)
		return err
	}
}

type OpCall func() error

func (call OpCall) Err() error {
	if call == nil {
		return nil
	}

	return call()
}

func (call OpCall) Must() {
	if err := call(); err != nil {
		panic(err)
	}
}

func (hs HostSlice) Len() int           { return len(hs) }
func (hs HostSlice) Less(i, j int) bool { return hs[i].ID() < hs[j].ID() }
func (hs HostSlice) Swap(i, j int)      { hs[i], hs[j] = hs[j], hs[i] }

func (hs HostSlice) Filter(f FilterFunc) HostSlice {
	out := make(HostSlice, 0, len(hs))
	for i, h := range hs {
		if f(i, h) {
			out = append(out, h)
		}
	}
	return out
}

func (hs HostSlice) Map(f func(i int, h host.Host) error) (err error) {
	for i, h := range hs {
		if err = f(i, h); err != nil {
			break
		}
	}

	return
}

func (hs HostSlice) Go(f func(i int, h host.Host) error) error {
	var g errgroup.Group
	for i, h := range hs {
		g.Go(func(i int, h host.Host) func() error {
			return func() (err error) {
				err = f(i, h)
				return
			}
		}(i, h))
	}
	return g.Wait()
}
