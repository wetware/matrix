package mx

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

func SelectIndex(i int) FilterFunc {
	return func(idx int, _ host.Host) bool {
		return i == idx
	}
}

func SelectIDs(ids ...peer.ID) FilterFunc {
	return func(_ int, h host.Host) (ok bool) {
		for _, id := range ids {
			if ok = h.ID() == id; ok {
				break
			}
		}

		return
	}
}
