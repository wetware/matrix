package namespace_test

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/stretchr/testify/require"
	"github.com/wetware/matrix/internal/testutil"
	"github.com/wetware/matrix/pkg/clock"
	"github.com/wetware/matrix/pkg/namespace"
)

const (
	name = "matrix.namespace.test"
	ttl  = time.Second
)

var t0 = time.Date(2021, 04, 9, 8, 0, 0, 0, time.UTC)

func TestProvider(t *testing.T) {
	t.Parallel()
	t.Helper()

	c := clock.New()
	c.Advance(t0)

	info := testutil.RandInfo()

	ns := namespace.New(c)
	got := ns.LoadOrCreate(name).
		Upsert(info, &discovery.Options{Ttl: ttl})

	c.Advance(t0.Add(ttl + c.Accuracy()))

	require.Equal(t, ttl, got)
	require.Eventually(t, func() bool {
		return len(ns.LoadOrCreate(name).Peers()) == 0
	}, time.Millisecond*100, time.Millisecond*10,
		"peer was not expired after %s", ttl)

}
