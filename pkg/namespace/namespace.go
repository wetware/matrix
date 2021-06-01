package namespace

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/wetware/matrix/pkg/netsim"
)

type Timer interface {
	After(d time.Duration, callback func()) (cancel func())
}

type Provider struct {
	mu sync.RWMutex
	m  map[string]netsim.Scope
	t  Timer
}

func New(t Timer) *Provider {
	return &Provider{m: make(map[string]netsim.Scope), t: t}
}

func (p *Provider) LoadOrCreate(name string) netsim.Scope {
	p.mu.RLock()
	ns, ok := p.Load(name)
	p.mu.RUnlock()

	if !ok { // slow path
		p.mu.Lock()
		defer p.mu.Unlock()

		// may have been added concurrently
		if ns, ok = p.m[name]; !ok {
			ns = p.namespace()
			p.m[name] = ns
		}
	}

	return ns
}

func (p *Provider) Load(name string) (ns netsim.Scope, ok bool) {
	p.mu.RLock()
	ns, ok = p.m[name]
	p.mu.RUnlock()
	return
}

type ns struct {
	mu sync.RWMutex
	rs map[peer.ID]*nsRecord
	t  Timer
}

func (p *Provider) namespace() *ns {
	return &ns{rs: make(map[peer.ID]*nsRecord), t: p.t}
}

func (ns *ns) Peers() netsim.InfoSlice {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	is := make(netsim.InfoSlice, 0, len(ns.rs))
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
	cancel func()
}

func record(info *peer.AddrInfo, cancel func()) *nsRecord {
	return &nsRecord{
		info:   info,
		cancel: cancel,
	}
}
