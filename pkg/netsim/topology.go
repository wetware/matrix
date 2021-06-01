package netsim

import (
	"context"
	"math/rand"
	"sort"
	"sync"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Topology selects peeers from the environment.
type Topology interface {
	SetDefaultOptions(*discovery.Options) error
	Select(context.Context, Namespace, *peer.AddrInfo, *discovery.Options) (InfoSlice, error)
}

type SelectAll struct{ defaultLoader }

// Implementations that embed SelectAll SHOULD call SelectAll.SetDefaultOptions before
// modifying the opts.
func (SelectAll) SetDefaultOptions(opts *discovery.Options) error {
	opts.Ttl = DefaultTTL
	opts.Other = make(map[interface{}]interface{})
	return nil
}

func (s SelectAll) Select(_ context.Context, ns Namespace, local *peer.AddrInfo, opts *discovery.Options) (InfoSlice, error) {
	return limit(opts, s.load(ns, local)), nil
}

type SelectRing struct{ SelectAll }

func (s SelectRing) Select(ctx context.Context, ns Namespace, local *peer.AddrInfo, opts *discovery.Options) (InfoSlice, error) {
	peers := s.load(ns, local)
	gt := peers.Filter(func(info *peer.AddrInfo) bool {
		return info.ID > local.ID
	})

	// largest peer?
	if len(gt) == 0 {
		return peers[0:1], nil // 'peers' is already sorted
	}

	return gt[0:1], nil
}

type SelectRandom struct {
	init sync.Once
	Src  rand.Source

	loader
	SelectAll
}

func (r *SelectRandom) Select(_ context.Context, ns Namespace, local *peer.AddrInfo, opts *discovery.Options) (InfoSlice, error) {
	r.init.Do(func() {
		if r.loader = (globalShuffleLoader{}); r.Src != nil {
			r.loader = &shuffleLoader{r: rand.New(r.Src)}
		}
	})

	return limit(opts, r.load(ns, local)), nil
}

func limit(opts *discovery.Options, as InfoSlice) InfoSlice {
	if opts.Limit == 0 || opts.Limit >= len(as) {
		return as
	}

	return as[:opts.Limit]
}

type loader interface {
	load(Namespace, *peer.AddrInfo) InfoSlice
}

// sortedLoader is embedded in various loaders/topologies (especially defaultLoader)
// in order to ensure reproducibility across runs.
func loadsort(ps interface{ Peers() InfoSlice }) InfoSlice {
	is := ps.Peers()
	sort.Sort(is)
	return is
}

// defaultLoader removes the local peer from the results
type defaultLoader struct{}

func (defaultLoader) load(ps interface{ Peers() InfoSlice }, local *peer.AddrInfo) InfoSlice {
	return loadsort(ps).
		Filter(func(info *peer.AddrInfo) bool { return info.ID != local.ID })
}

type globalShuffleLoader struct{}

func (globalShuffleLoader) load(ns Namespace, local *peer.AddrInfo) InfoSlice {
	return loadAndShuffle(ns, local, rand.Shuffle)
}

type shuffleLoader struct {
	mu sync.Mutex
	r  *rand.Rand
}

func (loader *shuffleLoader) load(ns Namespace, local *peer.AddrInfo) InfoSlice {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return loadAndShuffle(ns, local, loader.r.Shuffle)
}

func loadAndShuffle(ns Namespace, local *peer.AddrInfo, shuffle func(int, func(i, j int))) InfoSlice {
	as := defaultLoader{}.load(ns, local)
	shuffle(len(as), as.Swap)
	return as
}
