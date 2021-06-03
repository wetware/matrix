package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"github.com/libp2p/go-libp2p-core/network"
	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/netsim"
)

const (
	ns   = "matrix.example.basic"
	echo = "/echo"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim := mx.New(ctx)

	/*

	 Most Matrix functions have a corresponding Must* function
	 that panics instead of returning an error.  This provides
	 an (optional) way of reducing error-checking boilerplate.

	*/

	h0 := sim.MustHost(ctx)
	h0.SetStreamHandler(echo, handler)

	h1 := sim.MustHost(ctx)

	/*
	 Matrix provides the Operations API, which allows developers
	 to compose operations on collections of hosts.

	 Here, we're using a simple two-stage pipeline to announce
	 each peer to the namespace and connect them to each other.
	*/
	mx.Op(mx.Announce(sim, netsim.SelectAll{}, ns)).
		Then(mx.Discover(sim, netsim.SelectAll{}, ns)).
		MustArgs(ctx, h0, h1)

	s, err := h1.NewStream(ctx, h0.ID(), echo)
	maybeFatal(err)

	/*
	 Now we open a stream from h1 to h0, write some data, read
	 the response, and log it to stdout.
	*/

	var buf bytes.Buffer
	_, err = io.Copy(s, strings.NewReader("Hello, world!"))
	maybeFatal(err)
	maybeFatal(s.CloseWrite())

	_, err = io.Copy(&buf, s)
	maybeFatal(err)

	log.Println(buf.String())
}

func maybeFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handler(s network.Stream) {
	defer s.Close()

	_, err := io.Copy(s, s)
	maybeFatal(err)
}
