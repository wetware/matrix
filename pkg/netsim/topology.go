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
	Select(context.Context, Scope, *peer.AddrInfo, *discovery.Options) (InfoSlice, error)
}

type SelectAll struct{ defaultLoader }

// Implementations that embed SelectAll SHOULD call SelectAll.SetDefaultOptions before
// modifying the opts.
func (SelectAll) SetDefaultOptions(opts *discovery.Options) error {
	opts.Ttl = DefaultTTL
	opts.Other = make(map[interface{}]interface{})
	return nil
}

func (t SelectAll) Select(_ context.Context, s Scope, local *peer.AddrInfo, opts *discovery.Options) (InfoSlice, error) {
	return limit(opts, t.load(s, local)), nil
}

type SelectRing struct{ SelectAll }

func (t SelectRing) Select(ctx context.Context, s Scope, local *peer.AddrInfo, opts *discovery.Options) (InfoSlice, error) {
	peers := t.load(s, local)
	if len(peers) == 0 {
		return peers, nil
	}

	gt := peers.Filter(func(info *peer.AddrInfo) bool {
		return info.ID > local.ID
	})
	lt := peers.Filter(func(info *peer.AddrInfo) bool {
		return info.ID < local.ID
	})

	return append(gt, lt...), nil
}

type SelectRandom struct {
	init sync.Once
	Src  rand.Source

	loader
	SelectAll
}

func (t *SelectRandom) Select(_ context.Context, s Scope, local *peer.AddrInfo, opts *discovery.Options) (InfoSlice, error) {
	t.init.Do(func() {
		if t.loader = (globalShuffleLoader{}); t.Src != nil {
			t.loader = &shuffleLoader{r: rand.New(t.Src)}
		}
	})

	return limit(opts, t.load(s, local)), nil
}

func limit(opts *discovery.Options, as InfoSlice) InfoSlice {
	if opts.Limit == 0 || opts.Limit >= len(as) {
		return as
	}

	return as[:opts.Limit]
}

type loader interface {
	load(Scope, *peer.AddrInfo) InfoSlice
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

func (globalShuffleLoader) load(s Scope, local *peer.AddrInfo) InfoSlice {
	return loadAndShuffle(s, local, rand.Shuffle)
}

type shuffleLoader struct {
	mu sync.Mutex
	r  *rand.Rand
}

func (loader *shuffleLoader) load(s Scope, local *peer.AddrInfo) InfoSlice {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return loadAndShuffle(s, local, loader.r.Shuffle)
}

func loadAndShuffle(s Scope, local *peer.AddrInfo, shuffle func(int, func(i, j int))) InfoSlice {
	as := defaultLoader{}.load(s, local)
	shuffle(len(as), as.Swap)
	return as
}
