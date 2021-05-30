package main

import (
	"context"

	matrix "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/discover"
)

const ns = "casm.example.basic"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := matrix.New(ctx)

	/*

	 Most Matrix functions have a corresponding Must* function
	 that panics instead of returning an error.  This provides
	 an (optional) way of reducing error-checking boilerplate.

	*/
	h0 := env.MustHost(ctx)
	h1 := env.MustHost(ctx)

	env.Call(matrix.Announce(discover.SelectAll{}, ns)).
		Then(matrix.Discover(discover.SelectAll{}, ns)).
		Must(ctx, h0, h1)
}
