package main

import (
	"context"

	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/net"
)

const ns = "matrix.example.basic"

func opDiscover() mx.OpFunc {
	return mx.Announce(net.SelectAll{}, ns).
		Then(mx.Discover(net.SelectAll{}, ns))
}

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
	h1 := sim.MustHost(ctx)

	/*
	 Matrix provides the Operations API, which allows developers
	 to compose operations on collections of hosts.
	*/
	sim.Op(opDiscover()).
		Call(ctx, h0, h1).
		Must()
}
