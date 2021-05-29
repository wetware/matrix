package discover

import (
	"context"
	"errors"
	"math/rand"
	"sync"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Strategy selects peeers from the environment.
type Strategy interface {
	SetDefaultOptions(*discovery.Options) error
	Select(context.Context, Namespace, *discovery.Options) (InfoSlice, error)
}

type SelectAll struct{ nopOptionSetter }

func (s SelectAll) Select(_ context.Context, ns Namespace, opts *discovery.Options) (InfoSlice, error) {
	return limit(opts, ns.Peers()), nil
}

type SelectRing struct{ nopOptionSetter }

func (s SelectRing) Select(_ context.Context, ns Namespace, opts *discovery.Options) (InfoSlice, error) {
	id, ok := peerID(opts)
	if !ok {
		return nil, errors.New("ring topology requires option 'WithPeerID'")
	}

	var (
		is       = ns.Peers()
		neighbor *peer.AddrInfo
	)

	for i, info := range is {
		if id != info.ID {
			continue
		}

		// last peer?
		if i == len(is)-1 {
			neighbor = is[0] // wrap around to the beginning of the slice
			break
		}

		neighbor = is[i+1]
	}

	if neighbor == nil {
		return nil, errors.New("peer not in environment")
	}

	return InfoSlice{neighbor}, nil
}

type SelectRandom struct {
	init sync.Once
	Src  rand.Source

	loader
	nopOptionSetter
}

func (r *SelectRandom) Select(_ context.Context, ns Namespace, opts *discovery.Options) (InfoSlice, error) {
	r.init.Do(func() {
		if r.loader = (globalShuffleLoader{}); r.Src != nil {
			r.loader = &shuffleLoader{r: rand.New(r.Src)}
		}
	})

	return limit(opts, r.load(ns)), nil
}

func WithPeerID(id peer.ID) discovery.Option {
	return func(opts *discovery.Options) error {
		opts.Other[keyPeerID] = id
		return nil
	}
}

type key uint8

const (
	keyNamespace key = iota
	keyPeerID
)

func limit(opts *discovery.Options, as InfoSlice) InfoSlice {
	if opts.Limit == 0 || opts.Limit >= len(as) {
		return as
	}

	return as[:opts.Limit]
}

func peerID(opts *discovery.Options) (peer.ID, bool) {
	if v, ok := opts.Other[keyPeerID]; ok {
		return v.(peer.ID), true
	}

	return "", false
}

type nopOptionSetter struct{}

func (nopOptionSetter) SetDefaultOptions(*discovery.Options) error { return nil }

type loader interface {
	load(Namespace) InfoSlice
}

type globalShuffleLoader struct{}

func (globalShuffleLoader) load(ns Namespace) InfoSlice {
	return loadAndShuffle(ns, rand.Shuffle)
}

type shuffleLoader struct {
	mu sync.Mutex
	r  *rand.Rand
}

func (loader *shuffleLoader) load(ns Namespace) InfoSlice {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return loadAndShuffle(ns, loader.r.Shuffle)
}

func loadAndShuffle(ns Namespace, shuffle func(int, func(i, j int))) InfoSlice {
	as := ns.Peers()
	shuffle(len(as), as.Swap)
	return as
}
