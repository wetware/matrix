package netsim

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
)

var nopchan = make(chan peer.AddrInfo)

func init() { close(nopchan) }

const DefaultTTL = time.Hour * 24

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
	opts, err := d.options(ns, opt)
	if err != nil {
		return nil, err
	}

	n, ok := d.NS.Load(ns)
	if !ok {
		return nopchan, nil
	}

	is, err := d.Topo.Select(ctx, n, d.Info, opts)
	if err != nil {
		return nil, err
	}

	return infochan(is), nil
}

func (d *DiscoveryService) Advertise(ctx context.Context, ns string, opt ...discovery.Option) (time.Duration, error) {
	opts, err := d.options(ns, opt)
	if err != nil {
		return 0, err
	}

	return d.NS.LoadOrCreate(ns).Upsert(d.Info, opts), nil
}

func (d *DiscoveryService) options(ns string, opt []discovery.Option) (*discovery.Options, error) {
	d.init.Do(func() {
		d.validate = func(*discovery.Options) error { return nil }
		if v, ok := d.Topo.(validator); ok {
			d.validate = v.Validate
		}
	})

	var opts discovery.Options
	if err := d.Topo.SetDefaultOptions(&opts); err != nil {
		return nil, err
	}

	for _, option := range opt {
		if err := option(&opts); err != nil {
			return nil, err
		}
	}

	return &opts, d.validate(&opts)
}

func infochan(is InfoSlice) <-chan peer.AddrInfo {
	ch := make(chan peer.AddrInfo, len(is))
	defer close(ch)

	for _, info := range is {
		ch <- *info
	}

	return ch
}
