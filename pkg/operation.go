package mx

import (
	"github.com/libp2p/go-libp2p-core/host"
	"golang.org/x/sync/errgroup"
)

type OpFunc func(hs HostSlice) (HostSlice, error)

func (f OpFunc) Bind(fn OpFunc) Operation {
	return Op(func(hs HostSlice) (HostSlice, error) {
		hs, err := f(hs)
		if err != nil {
			return nil, err
		}

		return fn(hs)
	})
}

type Operation func(Operation) (OpFunc, Operation)

func (op Operation) Bind(fn func(OpFunc) Operation) Operation {
	return func(prev Operation) (OpFunc, Operation) {
		f, next := op(prev)
		return fn(f)(next)
	}
}

func (op Operation) Then(f OpFunc) Operation {
	return op.Bind(func(of OpFunc) Operation {
		return of.Bind(func(hs HostSlice) (HostSlice, error) {
			return f(hs)
		})
	})
}

func (op Operation) Eval(hs HostSlice) (out HostSlice, err error) {
	f, _ := op(Op(Just(nil)))
	return f(hs)
}

func (op Operation) Must(hs HostSlice) HostSlice {
	hs, err := op.Eval(hs)
	if err != nil {
		panic(err)
	}

	return hs
}

func (op Operation) Args(hs ...host.Host) (HostSlice, error) {
	return op.Eval(HostSlice(hs))
}

func (op Operation) MustArgs(hs ...host.Host) HostSlice {
	return op.Must(HostSlice(hs))
}

type HostSlice []host.Host

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
