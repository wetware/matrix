package env

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/discover"
)

type nsMap struct {
	mu sync.RWMutex
	m  map[string]discover.Namespace
	c  *Clock
}

func nsmap(c *clock.Clock) *nsMap {
	return &nsMap{m: make(map[string]discover.Namespace), c: (*Clock)(c)}
}

func (nm *nsMap) LoadOrCreate(name string) discover.Namespace {
	nm.mu.RLock()
	ns, ok := nm.Load(name)
	nm.mu.RUnlock()

	if !ok { // slow path
		nm.mu.Lock()
		defer nm.mu.Unlock()

		// may have been added concurrently
		if ns, ok = nm.m[name]; !ok {
			ns = namespace(nm.c)
			nm.m[name] = ns
		}
	}

	return ns
}

func (nm *nsMap) Load(name string) (ns discover.Namespace, ok bool) {
	nm.mu.RLock()
	ns, ok = nm.m[name]
	nm.mu.RUnlock()
	return
}

type ns struct {
	mu sync.RWMutex
	rs map[peer.ID]*nsRecord
	c  *Clock
}

func namespace(c *Clock) *ns { return &ns{rs: make(map[peer.ID]*nsRecord), c: c} }

func (ns *ns) Peers() discover.InfoSlice {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	is := make(discover.InfoSlice, 0, len(ns.rs))
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
	ns.rs[info.ID] = record(info, ns.c.After(opts.Ttl, func() {
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
