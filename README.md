# zzcache

[![CI](https://github.com/saromanov/zzcache/actions/workflows/go.yml/badge.svg)](https://github.com/saromanov/zzcache/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/saromanov/zzcache.svg)](https://pkg.go.dev/github.com/saromanov/zzcache)
[![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/zzcache)](https://goreportcard.com/report/github.com/saromanov/zzcache)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A sharded in-memory key/value cache for Go with per-entry TTL and pluggable stores.

Keys are distributed across shards by a CRC32 hash; each shard owns its own lock
and store instance, so operations on different keys do not contend. Expired
entries are evicted on read and optionally by a background janitor.

## Features

- Sharded: per-shard locks and per-shard store instances
- Per-entry TTL, zero duration means no expiration
- Expired entry eviction on read + optional background janitor
- Pluggable stores: default Go map, `radix` (prefix-friendly tree), custom via the `Store` interface
- Sentinel errors (`ErrNotFound`, `ErrExpired`, ...) for `errors.Is`
- `Count`, `Keys`, `Flush`, `Stats` (hits/misses)
- Race-detector clean; no external runtime dependencies except `go-radix`

## Installation

```
go get github.com/saromanov/zzcache
```

## Usage

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/saromanov/zzcache"
)

func main() {
	c, err := zzcache.New(512, "")
	if err != nil {
		panic(err)
	}

	if err := c.Set([]byte("key"), []byte("value"), time.Minute); err != nil {
		panic(err)
	}

	value, err := c.Get([]byte("key"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(value)) // value

	// Zero duration: the entry never expires.
	_ = c.Set([]byte("forever"), []byte("v"), 0)

	// Expired entries return ErrExpired and are evicted.
	_, err = c.Get([]byte("missing"))
	if errors.Is(err, zzcache.ErrNotFound) {
		fmt.Println("no such key")
	}
}
```

With the radix store (keeps keys in a prefix tree):

```go
c, err := zzcache.New(512, "radix")
```

## API

| Method | Description |
| --- | --- |
| `New(size uint32, storeType string) (*Cache, error)` | Create a cache with `size` shards; store type `""` (map) or `"radix"` |
| `Set(key, value []byte, d time.Duration) error` | Insert an entry; `d == 0` means no expiration |
| `Get(key []byte) ([]byte, error)` | Fetch an entry; `ErrNotFound` / `ErrExpired` |
| `Delete(key []byte) error` | Remove an entry; deleting a missing key is not an error |
| `Count() int` | Number of stored entries |
| `Keys() [][]byte` | Stored keys (order not defined) |
| `Flush()` | Remove all entries |
| `Stats() Stats` | Accumulated `Hits` / `Misses` |
| `StartCleanup(interval time.Duration) (stop func())` | Start a background janitor evicting expired entries |

Values are copied on `Set` and on `Get`: mutating a slice after `Set` or after
`Get` never corrupts the cache or other readers.

Keys are limited to 64 KB. A custom store can be built on the `Store` interface;
stores without `Range` support degrade `Count`/`Keys`/`Flush`/janitor to no-ops.

## Benchmarks

Go 1.25, linux/amd64, i7-9750H (512 shards, 10k preloaded keys, 128-byte values):

```
BenchmarkSet            	 3936043	       337.7 ns/op
BenchmarkSetParallel    	11405809	        99.6 ns/op
BenchmarkGet            	 4502907	       284.0 ns/op
BenchmarkGetParallel    	17167645	        75.5 ns/op
BenchmarkSetGetMixed    	 6725863	       186.5 ns/op
BenchmarkDelete         	 1798306	       576.0 ns/op
```

Reproduce with `go test -run XXX -bench . -benchtime=1s ./...`.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Pull requests are welcome.

## License

MIT — see [LICENSE](LICENSE).