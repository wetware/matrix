package mx

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"golang.org/x/sync/errgroup"
)

// Just discards the current selection and replaces it with hs.
func Just(hs Selection) SelectFunc {
	return func(context.Context, Selection) (Selection, error) {
		return hs, nil
	}
}

// Fail aborts the Op pipeline and returns the supplied error.
func Fail(err error) SelectFunc {
	return func(context.Context, Selection) (Selection, error) {
		return nil, err
	}
}

type SelectFunc func(ctx context.Context, hs Selection) (Selection, error)

func (f SelectFunc) Bind(fn SelectFunc) Op {
	return Select(func(ctx context.Context, hs Selection) (Selection, error) {
		hs, err := f(ctx, hs)
		if err != nil {
			return nil, err
		}

		return fn(ctx, hs)
	})
}

// NewPartition returns a partition of n subsets based on
// the index number of the current selection.
type Partition int

func (p Partition) Num() int { return int(p) }

// Get the partition by its index.
func (p Partition) Get(idx int) SelectFunc {
	return func(ctx context.Context, hs Selection) (Selection, error) {
		return hs.Filter(func(i int, _ host.Host) bool {
			return i%int(p) == idx
		}), nil
	}
}

// Selection is a set of hosts.
type Selection []host.Host

func (hs Selection) Len() int           { return len(hs) }
func (hs Selection) Less(i, j int) bool { return hs[i].ID() < hs[j].ID() }
func (hs Selection) Swap(i, j int)      { hs[i], hs[j] = hs[j], hs[i] }

func (hs Selection) Filter(f FilterFunc) Selection {
	out := make(Selection, 0, len(hs))
	for i, h := range hs {
		if f(i, h) {
			out = append(out, h)
		}
	}
	return out
}

func (hs Selection) Map(f func(i int, h host.Host) error) (err error) {
	for i, h := range hs {
		if err = f(i, h); err != nil {
			break
		}
	}

	return
}

func (hs Selection) Go(f func(i int, h host.Host) error) error {
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
