# Matrix

[![GoDoc](https://godoc.org/github.com/wetware/matrix?status.svg)](https://godoc.org/github.com/wetware/matrix)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)

In-process cluster simulation for libp2p.

Matrix is a library for **in-process** testing, benchmarking and simulating [libp2p](https://github.com/libp2p/go-libp2p) applicaitons in Go.  It is a simple alternative to [Testground](https://github.com/testground/testground) for cases where network IO is not needed or desired.

## Installation

```bash
go get -u github.com/wetware/matrix
```

## Motivation

### Testground

[Testground](https://github.com/testground/testground) is a platform for testing, benchmarking, and simulating distributed and p2p systems at scale. It's designed to be multi-lingual and runtime-agnostic, scaling gracefully from 2 to 10k instances, only when needed.

Testground is ideal for complex workloads, such as:

- Compatibility testing between different versions of a P2P application
- Verifying interoperability between different language implementations of a library
- Simulations that require real network IO.

The price to pay for Testground's minimal assumptions and maximal realism come with some limitations and additional complexities.  Test plans:

- require configuration management (`.env.toml`, `manifest.toml`)
- require a specialized testing environment (Tesground daemon, Redis, Docker, etc.)
- write data to the local filesystem
- cannot be embedded in unit tests (no support for `go test`)

Taken together, these characteristics can make Testground overkill for simple unit-tests and benchmarks, especially if these are integrated into your local development workflow (or if you expect contributors to run your tests).  In this common case, an in-process simulation suffices.

### Matrix

Matrix is a library for writing unit-test, benchmarks and exploratory simulations for libp2p.  Contrary to Testground, everything happens in a single process, without using the network.

#### Goals:
- Drop-in compatibility with PL stack
- In-process.  No external services / sidecar processes / environmental dependencies
- Support for benchmarking, unit-testing and runtime analysis
- Traffic shaping
- Stats collecting / offline analysis

#### Non-Goals
- Simulated time
- Multilingual Support

Matrix runs your code using actual libp2p code.  It configures [hosts](https://pkg.go.dev/github.com/libp2p/go-libp2p-core/host#Host) to use an [in-process transport](https://godoc.org/github.com/lthibault/go-libp2p-inproc-transport), which allows them to communicate without the network.  **Everything else is exactly the same.**

Additionally, Matrix provides some utilities to facilitate test setup.  For example, the [`discover`](pkg.go.dev/github.com/wetware/matrix/pkg/discover) package provides a specialized [`discovery.Discovery`](https://pkg.go.dev/github.com/libp2p/go-libp2p-core/discovery#Discovery) implementation that allows you to arrange hosts in a specific topology (e.g. a ring).

But there's more!  Just like Testground, Matrix provides support for sophisticated traffic shaping.

## Usage

The following example can be found under `examples/basic`.  See the examples directory for more.

```go
import (
    "context"

    mx "github.com/wetware/matrix/pkg"
    "github.com/wetware/matrix/pkg/net"
)

const ns   = "matrix.test"


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

    Here, we're using a simple two-stage pipeline to announce
    each peer to the namespace and connect them to each other.
*/
sim.Op(mx.Announce(net.SelectAll{}, ns)).
    Then(mx.Discover(net.SelectAll{}, ns)).
    Call(ctx, h0, h1).
    Must()

/*
    h0 and h1 are now connected to each other!
*/
```

## Team

### Core Team

- [@lthibault](https://github.com/lthibault) ★
- [@aratz-lasa](https://github.com/aratz-lasa)

★ Project Lead

## License

Dual-licensed: [MIT](https://github.com/testground/testground/blob/master/LICENSE-MIT), [Apache Software License v2](https://github.com/testground/testground/blob/master/LICENSE-APACHE), by way of the [Permissive License Stack](https://protocol.ai/blog/announcing-the-permissive-license-stack/).
