package mx_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/libp2p/go-libp2p-core/network"
	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/netsim"
)

func ExampleSimulation() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		ns   = "matrix.test"
		echo = "/echo"
	)

	var (
		sim = mx.New(ctx)

		h0 = sim.MustHost(ctx)
		h1 = sim.MustHost(ctx)
	)

	h0.SetStreamHandler(echo, func(s network.Stream) {
		defer s.Close()

		if _, err := io.Copy(s, s); err != nil {
			panic(err)
		}
	})

	mx.Topology(sim, netsim.SelectRing{}, ns).
		MustArgs(ctx, h0, h1)

	s, err := h1.NewStream(ctx, h0.ID(), echo)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	_, err = io.Copy(s, strings.NewReader("Hello, world!"))
	if err != nil {
		panic(err)
	}
	s.CloseWrite()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, s)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
	// Output: Hello, world!
}
