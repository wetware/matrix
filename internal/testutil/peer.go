package testutil

import (
	"fmt"
	"math/rand"

	"github.com/libp2p/go-libp2p-core/peer"
	inproc "github.com/lthibault/go-libp2p-inproc-transport"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
)

func RandInfo() *peer.AddrInfo {
	id := NewID()
	return &peer.AddrInfo{
		ID:    id,
		Addrs: []multiaddr.Multiaddr{Multiaddr(id)},
	}
}

func Multiaddr(id peer.ID) multiaddr.Multiaddr {
	ma, err := inproc.ResolveString("/inproc/~")
	if err != nil {
		panic(err)
	}

	return ma.Encapsulate(multiaddr.StringCast(fmt.Sprintf("/p2p/%s", id)))
}

func NewID() peer.ID {
	return newID(randStr(5))
}

func randStr(n int) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	b := make([]rune, n)
	for i := range b {
		b[i] = rune(alphabet[rand.Intn(len(alphabet))])
	}

	return string(b)
}

func hash(b []byte) []byte {
	h, _ := multihash.Sum(b, multihash.SHA2_256, -1)
	return []byte(h)
}

func newID(s string) peer.ID {
	id, err := peer.Decode(base58.Encode(hash([]byte(s))))
	if err != nil {
		panic(err)
	}

	return id
}
