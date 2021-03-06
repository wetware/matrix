package mx

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

type FilterFunc func(int, host.Host) bool

// Nop returns the selection unchanged.
func Nop() Op {
	return Select(func(_ context.Context, hs Selection) (Selection, error) {
		return hs, nil
	})
}

func Select(f SelectFunc) Op {
	return func(op Op) (SelectFunc, Op) {
		return f, op
	}
}

// Map applies f to each item in the current selection.
func Map(f MapFunc) Op {
	return mapper(func(hs Selection, hf func(MapFunc) func(int, host.Host) error) error {
		return hs.Map(hf(f))
	})
}

// Go applies f to each item in the current selection concurrently.
func Go(f MapFunc) Op {
	return mapper(func(hs Selection, hf func(MapFunc) func(int, host.Host) error) error {
		return hs.Go(hf(f))
	})
}

// Filter returns a new selection that contains the elements of the
// current selection for which f(element) == true.
func Filter(f FilterFunc) Op {
	return Select(func(ctx context.Context, hs Selection) (Selection, error) {
		return hs.Filter(f), nil
	})
}

func connect(ctx context.Context, h host.Host, info peer.AddrInfo) func() error {
	return func() error {
		return h.Connect(ctx, info)
	}
}

func mapper(f func(hs Selection, hf func(MapFunc) func(int, host.Host) error) error) Op {
	return Select(func(ctx context.Context, hs Selection) (Selection, error) {
		err := f(hs, func(mf MapFunc) func(i int, h host.Host) error {
			return func(i int, h host.Host) error {
				return mf(ctx, i, h)
			}
		})

		if err != nil {
			return nil, err
		}

		return hs, nil
	})
}

type Op func(Op) (SelectFunc, Op)

func (op Op) Bind(fn func(SelectFunc) Op) Op {
	return func(prev Op) (SelectFunc, Op) {
		f, next := op(prev)
		return fn(f)(next)
	}
}

func (op Op) Then(f SelectFunc) Op {
	return op.Bind(func(of SelectFunc) Op {
		return of.Bind(func(ctx context.Context, hs Selection) (Selection, error) {
			return f(ctx, hs)
		})
	})
}

func (op Op) Map(f MapFunc) Op {
	return op.Bind(func(sf SelectFunc) Op {
		return sf.Bind(func(ctx context.Context, hs Selection) (Selection, error) {
			return hs, hs.Map(func(i int, h host.Host) error {
				return f(ctx, i, h)
			})
		})
	})
}

func (op Op) Go(f MapFunc) Op {
	return op.Bind(func(sf SelectFunc) Op {
		return sf.Bind(func(ctx context.Context, hs Selection) (Selection, error) {
			return hs, hs.Go(func(i int, h host.Host) error {
				return f(ctx, i, h)
			})
		})
	})
}

func (op Op) Eval(ctx context.Context, hs Selection) (out Selection, err error) {
	f, _ := op(Select(Just(nil)))
	return f(ctx, hs)
}

func (op Op) Err(ctx context.Context, hs Selection) (err error) {
	_, err = op.Eval(ctx, hs)
	return
}

func (op Op) Must(ctx context.Context, hs Selection) Selection {
	hs, err := op.Eval(ctx, hs)
	if err != nil {
		panic(err)
	}

	return hs
}

func (op Op) EvalArgs(ctx context.Context, hs ...host.Host) (Selection, error) {
	return op.Eval(ctx, hs)
}

func (op Op) ErrArgs(ctx context.Context, hs ...host.Host) error {
	return op.Err(ctx, hs)
}

func (op Op) MustArgs(ctx context.Context, hs ...host.Host) Selection {
	return op.Must(ctx, hs)
}

type MapFunc func(ctx context.Context, i int, h host.Host) error

func (f MapFunc) Then(fn MapFunc) MapFunc {
	return func(ctx context.Context, i int, h host.Host) error {
		if err := f(ctx, i, h); err != nil {
			return err
		}

		return fn(ctx, i, h)
	}
}
