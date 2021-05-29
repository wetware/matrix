package discover

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/peer"
)

const DefaultTTL = time.Hour * 8766

var (
	_ discovery.Discovery = (*discoverer)(nil)

	nopchan = make(chan peer.AddrInfo)
)

func init() { close(nopchan) }

type Namespace interface {
	Peers() InfoSlice
	Upsert(*peer.AddrInfo, *discovery.Options) time.Duration
}

type NamespaceProvider interface {
	Load(ns string) (Namespace, bool)
	LoadOrCreate(ns string) Namespace
}

type discoverer struct {
	ns       NamespaceProvider
	info     *peer.AddrInfo
	s        Strategy
	validate func(*discovery.Options) error
}

func New(ns NamespaceProvider, s Strategy, info *peer.AddrInfo) discovery.Discovery {
	d := discoverer{
		ns:       ns,
		info:     info,
		s:        s,
		validate: func(*discovery.Options) error { return nil },
	}

	if v, ok := d.s.(validator); ok {
		d.validate = v.Validate
	}

	return d
}

type validator interface {
	Validate(*discovery.Options) error
}

func (d discoverer) FindPeers(ctx context.Context, ns string, opt ...discovery.Option) (<-chan peer.AddrInfo, error) {
	n, ok := d.ns.Load(ns)
	if !ok {
		return nopchan, nil
	}

	opts, err := d.options(ns, opt)
	if err != nil {
		return nil, err
	}

	as, err := d.s.Select(ctx, n, opts)
	return infochan(as), err
}

func (d discoverer) Advertise(ctx context.Context, ns string, opt ...discovery.Option) (time.Duration, error) {
	opts, err := options(opt)
	if err != nil {
		return 0, err
	}

	return d.ns.LoadOrCreate(ns).Upsert(d.info, opts), nil
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

func (d discoverer) options(ns string, opt []discovery.Option) (*discovery.Options, error) {
	opts := newOptions()
	if err := d.s.SetDefaultOptions(opts); err != nil {
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

type InfoSlice []*peer.AddrInfo

func (is InfoSlice) Len() int           { return len(is) }
func (is InfoSlice) Less(i, j int) bool { return is[i].ID < is[j].ID }
func (is InfoSlice) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
