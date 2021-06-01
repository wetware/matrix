package mx

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/require"

	"github.com/wetware/matrix/internal/testutil"
)

func TestSimulation(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.SkipNow()
}

func TestNewDiscovery(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim := New(ctx)

	h := testutil.NewHost(ctrl)

	svc := sim.NewDiscovery(h, nil)
	require.NotNil(t, svc)
	require.NotNil(t, svc.Topo)
	require.Equal(t, svc.Info, host.InfoFromHost(h))
}
