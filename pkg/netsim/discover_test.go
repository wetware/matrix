package netsim

// func TestDiscovery(t *testing.T) {
// 	t.Parallel()
// 	t.Helper()

// 	t.Run("DefaultOptionErrorFails", func(t *testing.T) {
// 		t.Parallel()

// 		d := discoveryService{topo: failDefaultOptions{}}
// 		peers, err := d.FindPeers(context.Background(), "")
// 		require.EqualError(t, err, "test")
// 		require.Nil(t, peers)
// 	})

// 	t.Run("BadOptionFails", func(t *testing.T) {
// 		t.Parallel()

// 		var d discoveryService
// 		peers, err := d.FindPeers(context.Background(), "",
// 			func(*discovery.Options) error { return errors.New("test") })
// 		require.EqualError(t, err, "test")
// 		require.Nil(t, peers)
// 	})

// 	t.Run("ValidationErrorFails", func(t *testing.T) {
// 		t.Parallel()

// 		d := discoveryService{topo: failValidaton{}}
// 		peers, err := d.FindPeers(context.Background(), "")
// 		require.EqualError(t, err, "test")
// 		require.Nil(t, peers)
// 	})

// 	t.Run("Succeed", func(t *testing.T) {
// 		t.Parallel()

// 		const n = 10
// 		newTestEnv(n)

// 		d := discoveryService{
// 			Env:      newTestEnv(n),
// 			Strategy: netsim.SelectAll{},
// 		}

// 		peers, err := d.FindPeers(context.Background(), "")
// 		require.NoError(t, err)
// 		require.Len(t, peers, n)
// 	})
// }

// func TestStrategy(t *testing.T) {
// 	t.Parallel()
// 	t.Helper()

// 	const n = 10
// 	env := newTestEnv(n)

// 	t.Run("SelectAll", func(t *testing.T) {
// 		t.Parallel()
// 		t.Helper()

// 		t.Run("Limit", func(t *testing.T) {
// 			t.Parallel()

// 			const limit = 5

// 			var s netsim.SelectAll

// 			as, err := runStrategy(env, s, discovery.Limit(limit))
// 			require.NoError(t, err)

// 			require.Subset(t, env.List(), as)
// 			require.Len(t, as, limit)

// 			peers := env.List()
// 			sort.Sort(peers)
// 			assert.Equal(t, peers[:limit], as)
// 		})
// 	})

// 	t.Run("SelectRandom", func(t *testing.T) {
// 		t.Parallel()
// 		t.Helper()

// 		t.Run("GlobalSource", func(t *testing.T) {
// 			t.Parallel()

// 			var s netsim.SelectRandom
// 			as, err := runStrategy(env, &s)
// 			require.NoError(t, err)

// 			require.ElementsMatch(t, env.List(), as)
// 		})

// 		t.Run("Reproducible", func(t *testing.T) {
// 			t.Parallel()

// 			as0, err := runStrategy(env, &netsim.SelectRandom{
// 				Src: rand.NewSource(42),
// 			})
// 			require.NoError(t, err)

// 			as1, err := runStrategy(env, &netsim.SelectRandom{
// 				Src: rand.NewSource(42),
// 			})
// 			require.NoError(t, err)

// 			require.ElementsMatch(t, as0, as1)
// 			require.Equal(t, as0, as1)
// 		})
// 	})

// 	t.Run("SelectRing", func(t *testing.T) {
// 		t.Parallel()
// 		t.Helper()

// 		t.Run("MissingPeerIDFails", func(t *testing.T) {
// 			t.Parallel()

// 			var s netsim.SelectRing
// 			as, err := runStrategy(env, s)
// 			require.Error(t, err)
// 			require.Nil(t, as)
// 		})

// 		t.Run("PeerNotInEnvironmentFails", func(t *testing.T) {
// 			t.Parallel()

// 			var s netsim.SelectRing
// 			as, err := runStrategy(env, s, net.WithPeerID(randID()))
// 			require.Error(t, err)
// 			require.Nil(t, as)
// 		})

// 		t.Run("Succeeds", func(t *testing.T) {
// 			t.Parallel()

// 			var (
// 				peers     = env.List()
// 				neighbors = make(inproc.AddrSlice, len(peers))
// 			)

// 			sort.Sort(peers)

// 			var g errgroup.Group
// 			for i, a := range peers {
// 				g.Go(func(i int, a multiaddr.Multiaddr) func() error {
// 					return func() (err error) {
// 						var (
// 							info *peer.AddrInfo
// 							s    netsim.SelectRing
// 							as   inproc.AddrSlice
// 						)
// 						if info, err = peer.AddrInfoFromP2pAddr(a); err != nil {
// 							return
// 						}

// 						if as, err = runStrategy(env, s, net.WithPeerID(info.ID)); err == nil {
// 							neighbors[i] = as[0]
// 						}

// 						return
// 					}
// 				}(i, a))
// 			}

// 			require.NoError(t, g.Wait())
// 			require.Equal(t, peers,
// 				// same array, except that tail is appended to head.
// 				append(neighbors[n-1:], neighbors[:n-1]...))
// 		})
// 	})
// }

// func runStrategy(env net.PeerListProvider, s net.Strategy, opt ...discovery.Option) (inproc.AddrSlice, error) {
// 	opts := newOption()
// 	if err := s.SetDefaultOptions(opts); err != nil {
// 		return nil, err
// 	}

// 	for _, option := range opt {
// 		if err := option(opts); err != nil {
// 			return nil, err
// 		}
// 	}

// 	// validate?
// 	if v, ok := s.(interface {
// 		Validate(*discovery.Options) error
// 	}); ok {
// 		if err := v.Validate(opts); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return s.Select(context.Background(), opts, env)
// }

// func newOption() *discovery.Options {
// 	return &discovery.Options{Other: make(map[interface{}]interface{})}
// }

// func newTestEnv(n int) inproc.Env {
// 	env := inproc.NewEnv()
// 	for i := 0; i < n; i++ {
// 		if !env.Bind(newAddr(), new(inproc.Transport)) {
// 			panic("failed to bind")
// 		}
// 	}
// 	return env
// }

// func newAddr() multiaddr.Multiaddr {
// 	ma, err := inproc.ResolveString("/inproc/~")
// 	if err != nil {
// 		panic(err)
// 	}

// 	return ma.Encapsulate(multiaddr.StringCast(fmt.Sprintf("/p2p/%s", randID())))
// }

// func randID() peer.ID {
// 	return newID(randStr(5))
// }

// func randStr(n int) string {
// 	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = rune(alphabet[rand.Intn(len(alphabet))])
// 	}

// 	return string(b)
// }

// func hash(b []byte) []byte {
// 	h, _ := multihash.Sum(b, multihash.SHA2_256, -1)
// 	return []byte(h)
// }

// func newID(s string) peer.ID {
// 	id, err := peer.Decode(base58.Encode(hash([]byte(s))))
// 	if err != nil {
// 		panic(err)
// 	}

// 	return id
// }

// type failValidaton struct{ netsim.SelectAll }

// func (failValidaton) Validate(*discovery.Options) error {
// 	return errors.New("test")
// }

// type failDefaultOptions struct{ netsim.SelectAll }

// func (failDefaultOptions) SetDefaultOptions(*discovery.Options) error {
// 	return errors.New("test")
// }
