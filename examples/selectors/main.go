package main

import (
	"context"
	"fmt"

	mx "github.com/wetware/matrix/pkg"
	"github.com/wetware/matrix/pkg/netsim"
)

const (
	n    = 10
	ns   = "matrix.example.basic"
	echo = "/echo"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		sim = mx.New(ctx)
		hs  = sim.MustHostSet(ctx, n) // create slice of n hosts
		p   = mx.Partition(2)         // partition of size 2
	)

	p0 := mx.Topology(sim, netsim.SelectRing{}, ns).
		Then(p.Get(0)). // select the first partition
		Must(ctx, hs)

	fmt.Printf("partition of %d hosts:\n\n", len(p0))
	for i, h := range p0 {
		fmt.Printf("%d\t%s\n", i, h.ID())
	}
}
