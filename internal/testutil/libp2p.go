package testutil

import (
	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	ma "github.com/multiformats/go-multiaddr"
	mock_libp2p "github.com/wetware/matrix/internal/mock/libp2p"
)

func NewHost(ctrl *gomock.Controller) *mock_libp2p.MockHost {
	h := mock_libp2p.NewMockHost(ctrl)
	h.EXPECT().ID().Return(PeerID()).AnyTimes()
	h.EXPECT().Addrs().Return(addrsFor(h.ID())).AnyTimes()
	h.EXPECT().Network().Return(networkFor(ctrl, h)).AnyTimes()
	return h
}

// PeerID returns a random peer.ID for testing
func PeerID() peer.ID {
	return newID(randStr(5))
}

func addrsFor(id peer.ID) []ma.Multiaddr {
	addr, _ := inproc.ResolveString("/inproc/~")
	return []ma.Multiaddr{addr}
}

func networkFor(ctrl *gomock.Controller, h host.Host) *mock_libp2p.MockNetwork {
	n := mock_libp2p.NewMockNetwork(ctrl)
	n.EXPECT().LocalPeer().Return(h.ID()).AnyTimes()
	n.EXPECT().ListenAddresses().Return([]ma.Multiaddr{ma.StringCast("/inproc/~")}).AnyTimes()
	n.EXPECT().InterfaceListenAddresses().Return(h.Addrs(), nil).AnyTimes()
	return n
}
