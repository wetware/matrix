package namespace_test

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/internal/testutil"
	"github.com/wetware/matrix/pkg/namespace"
)

const (
	name = "matrix.namespace.test"
	ttl  = time.Second
)

func TestProvider(t *testing.T) {
	t.Parallel()
	t.Helper()

	c := make(chan struct{})

	info := testutil.RandInfo()

	ns := namespace.New(mockTimer(c))
	got := ns.LoadOrCreate(name).
		Upsert(info, &discovery.Options{Ttl: ttl})

	close(c) // signal expiration

	require.Equal(t, ttl, got)
	require.Eventually(t, func() bool {
		return len(ns.LoadOrCreate(name).Peers()) == 0
	}, time.Millisecond*100, time.Millisecond*10,
		"peer was not expired after %s", ttl)

}

type mockTimer <-chan struct{}

func (t mockTimer) After(_ time.Duration, callback func()) func() {
	cancel := make(chan struct{})
	go func() {
		select {
		case <-t:
			callback()
		case <-cancel:
		}
	}()
	return func() { close(cancel) }
}
