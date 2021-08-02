//go:generate mockgen -destination ../../internal/mock/pkg/netsim/matrix.go github.com/wetware/matrix/pkg/netsim Scope,NamespaceProvider

package netsim

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/config"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
)

type Scope interface {
	Peers() InfoSlice
	Upsert(*peer.AddrInfo, *discovery.Options) time.Duration
}

type NamespaceProvider interface {
	LoadOrCreate(string) Scope
	Load(string) (Scope, bool)
}

type InfoSlice []*peer.AddrInfo

func (is InfoSlice) Len() int           { return len(is) }
func (is InfoSlice) Less(i, j int) bool { return is[i].ID < is[j].ID }
func (is InfoSlice) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

func (is InfoSlice) Filter(f func(info *peer.AddrInfo) bool) InfoSlice {
	filt := make(InfoSlice, 0, len(is))
	for _, info := range is {
		if f(info) {
			filt = append(filt, info)
		}
	}
	return filt
}

type HostFactory struct {
	init sync.Once
	inproc.Env
}

// NewHost assembles and creates a new libp2p host that uses the
// simulation's network.
//
// Env configures hosts to use an in-process network and therefore
// overrides the following options:
//
// - libp2p.Transport
// - libp2p.NoTransports
// - libp2p.ListenAddr
// - libp2p.ListenAddrStrings
// - libp2p.NoListenAddrs
//
// Users SHOULD NOT pass any of the above options to NewHost.
func (f *HostFactory) NewHost(ctx context.Context, opt []config.Option) (host.Host, error) {
	f.init.Do(func() { f.Env = inproc.NewEnv() })
	return libp2p.New(ctx, f.hostopt(opt)...)
}

func (f *HostFactory) hostopt(opt []libp2p.Option) []config.Option {
	return append([]libp2p.Option{
		libp2p.NoListenAddrs,
		libp2p.ListenAddrStrings("/inproc/~"),
		libp2p.NoTransports,
		libp2p.Transport(inproc.New(inproc.WithEnv(f))),
	}, opt...)
}
