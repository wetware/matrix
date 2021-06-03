package mx_test

// func TestPartition(t *testing.T) {
// 	t.Parallel()
// 	t.Helper()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	sim := mx.New(ctx,
// 		mx.WithClock(testutil.NewClock(ctrl, 0, nil)),
// 		mx.WithHostFactory(testutil.NewHostFactory(ctrl)))

// 	hs := sim.MustHostSet(ctx, n)

// 	t.Run("Num", func(t *testing.T) {
// 		t.Parallel()

// 		for i := 0; i < n; i++ {
// 			require.Equal(t, i, mx.Partition(i).Num())
// 		}
// 	})

// 	t.Run("Singletons", func(t *testing.T) {
// 		t.Parallel()

// 		p := mx.Partition(n)

// 		for i := 0; i < p.Num(); i++ {
// 			res := sim.Op(p.At(i)).Must(ctx, hs...)
// 			require.Len(t, res, 1)
// 			require.Equal(t, hs[i], res[0])
// 		}
// 	})

// 	t.Run("Bipartite", func(t *testing.T) {
// 		t.Parallel()

// 		p := mx.Partition(2)

// 		var ps []mx.HostSlice
// 		for i := 0; i < 2; i++ {
// 			res := sim.Op(p.At(i)).
// 				Then(mx.Select(func(ctx context.Context, sim mx.Simulation, hs mx.HostSlice) (mx.HostSlice, error) {
// 					ps = append(ps, hs)
// 					return hs, nil
// 				})).
// 				Must(ctx, hs...)

// 			require.Len(t, res, n/2)
// 		}

// 		require.Len(t, ps, 2)

// 		// TODO:  check that the partitions don't overlap
// 	})
// }
