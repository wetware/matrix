package env_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wetware/matrix/internal/testutil"
	"github.com/wetware/matrix/pkg/discover"
	"github.com/wetware/matrix/pkg/env"
)

func TestEnv(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := env.New(ctx)

	require.NotNil(t, e.Clock())
	require.NotNil(t, e.Network())
	require.NotNil(t, e.Process())

	t.Run("NewDiscovery", func(t *testing.T) {
		t.Parallel()

		h := testutil.NewHost(ctrl)

		d := e.NewDiscovery(*host.InfoFromHost(h), discover.SelectAll{})

		t.Run("Announce", func(t *testing.T) {
			ttl, err := d.Advertise(ctx, "test")
			require.NoError(t, err)
			assert.Equal(t, discover.DefaultTTL, ttl)
		})

		t.Run("FindPeers", func(t *testing.T) {
			ch, err := d.FindPeers(ctx, "test")
			require.NoError(t, err)
			require.Len(t, ch, 1)
			assert.Equal(t, h.ID(), (<-ch).ID)
		})
	})
}
