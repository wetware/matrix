package netsim

import (
	"sort"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/wetware/matrix/pkg/clock"
)

type nsMap struct {
	mu sync.RWMutex
	m  map[string]Namespace
	t  Timer
}

func nsmap(t Timer) *nsMap {
	return &nsMap{m: make(map[string]Namespace), t: t}
}

func (nm *nsMap) LoadOrCreate(name string) Namespace {
	nm.mu.RLock()
	ns, ok := nm.Load(name)
	nm.mu.RUnlock()

	if !ok { // slow path
		nm.mu.Lock()
		defer nm.mu.Unlock()

		// may have been added concurrently
		if ns, ok = nm.m[name]; !ok {
			ns = namespace(nm.t)
			nm.m[name] = ns
		}
	}

	return ns
}

func (nm *nsMap) Load(name string) (ns Namespace, ok bool) {
	nm.mu.RLock()
	ns, ok = nm.m[name]
	nm.mu.RUnlock()
	return
}

type ns struct {
	mu sync.RWMutex
	rs map[peer.ID]*nsRecord
	t  Timer
}

func namespace(t Timer) *ns {
	return &ns{rs: make(map[peer.ID]*nsRecord), t: t}
}

func (ns *ns) Peers() InfoSlice {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	is := make(InfoSlice, 0, len(ns.rs))
	defer sort.Sort(is) // ensure reproducibility

	for _, rec := range ns.rs {
		is = append(is, rec.info)
	}

	return is
}

func (ns *ns) Upsert(info *peer.AddrInfo, opts *discovery.Options) time.Duration {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if rec, ok := ns.rs[info.ID]; ok {
		rec.cancel()
	}

	return ns.insert(info, opts)
}

func (ns *ns) insert(info *peer.AddrInfo, opts *discovery.Options) time.Duration {
	ns.rs[info.ID] = record(info, ns.t.After(opts.Ttl, func() {
		ns.mu.Lock()
		delete(ns.rs, info.ID)
		ns.mu.Unlock()
	}))

	return opts.Ttl
}

type nsRecord struct {
	info   *peer.AddrInfo
	cancel clock.CancelFunc
}

func record(info *peer.AddrInfo, cancel clock.CancelFunc) *nsRecord {
	return &nsRecord{
		info:   info,
		cancel: cancel,
	}
}
