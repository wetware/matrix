package net

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
)

var nopchan = make(chan peer.AddrInfo)

func init() { close(nopchan) }

type Namespace interface {
	Peers() InfoSlice
	Upsert(*peer.AddrInfo, *discovery.Options) time.Duration
}

type NamespaceProvider interface {
	Load(ns string) (Namespace, bool)
	LoadOrCreate(ns string) Namespace
}

type InfoSlice []*peer.AddrInfo

func (is InfoSlice) Len() int           { return len(is) }
func (is InfoSlice) Less(i, j int) bool { return is[i].ID < is[j].ID }
func (is InfoSlice) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

const DefaultTTL = time.Hour * 8766

type DiscoveryService struct {
	NS   NamespaceProvider
	Info *peer.AddrInfo
	Topo Topology

	init     sync.Once
	validate func(*discovery.Options) error
}

type validator interface {
	Validate(*discovery.Options) error
}

func (d *DiscoveryService) FindPeers(ctx context.Context, ns string, opt ...discovery.Option) (<-chan peer.AddrInfo, error) {
	d.init.Do(func() {
		if d.Topo == nil {
			d.Topo = SelectAll{}
		}

		d.validate = func(o *discovery.Options) error { return nil }
		if v, ok := d.Topo.(validator); ok {
			d.validate = v.Validate
		}
	})

	n, ok := d.NS.Load(ns)
	if !ok {
		return nopchan, nil
	}

	opts, err := d.options(ns, opt)
	if err != nil {
		return nil, err
	}

	as, err := d.Topo.Select(ctx, n, opts)
	return infochan(as), err
}

func (d *DiscoveryService) Advertise(ctx context.Context, ns string, opt ...discovery.Option) (time.Duration, error) {
	opts, err := options(opt)
	if err != nil {
		return 0, err
	}

	return d.NS.LoadOrCreate(ns).Upsert(d.Info, opts), nil
}

func options(opt []discovery.Option) (*discovery.Options, error) {
	opts := &discovery.Options{Ttl: DefaultTTL}
	for _, option := range opt {
		if err := option(opts); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func (d *DiscoveryService) options(ns string, opt []discovery.Option) (*discovery.Options, error) {
	opts := newOptions()
	if err := d.Topo.SetDefaultOptions(opts); err != nil {
		return nil, err
	}

	for _, option := range opt {
		if err := option(opts); err != nil {
			return nil, err
		}
	}

	return opts, d.validate(opts)
}

func infochan(is InfoSlice) <-chan peer.AddrInfo {
	ch := make(chan peer.AddrInfo, len(is))
	defer close(ch)

	for _, info := range is {
		ch <- *info
	}

	return ch
}

func newOptions() *discovery.Options {
	return &discovery.Options{Other: make(map[interface{}]interface{})}
}
