# cloudranger

`cloudranger` is a Go library designed to identify cloud provider information from IP addresses, with current support for AWS and GCP.

It functions without any external runtime dependencies, as IP range data is stored internally. Meant for high throughput, low-latency environments, `cloudranger` also focuses on rapid startup, loading in under 4ms. You can verify this on your system by running `make bench` and checking the `BenchmarkNew` results.

New releases are automatically created in response to updates in cloud providers' IP range information. This process, facilitated through GitHub Actions, is executed weekly to ensure the library remains up-to-date.

Inspiration for this library came from https://github.com/kubernetes/registry.k8s.io which contains a similar library used by the Kubernetes OCI registry to redirect requests to the appropriate cloud provider. `coderanger` was created as a means to have a standalone library that can be used in other projects and with greater control for our own use cases. That project also uses its own trie implementation whereas this library depends on `github.com/infobloxopen/go-trees`. I have not benchmarked the two implementations, but I suspect they are comparable.

## Usage

```sh
go get github.com/planetscale/cloudranger
```

```go
package main

import (
	"fmt"

	"github.com/planetscale/cloudranger"
)

func main() {
	ranger := cloudranger.New()
	ipinfo, found := ranger.GetIP("3.5.140.101")
	if found {
		fmt.Printf("cloud: %s, region: %s\n", ipinfo.Cloud(), ipinfo.Region())
	}
}
```

## Testing and Benchmarks

```sh
make lint
make test
make bench
```

```
goos: linux
goarch: amd64
pkg: github.com/planetscale/cloudranger
cpu: AMD EPYC 7B12

BenchmarkNew-16              396           3153591 ns/op         1084213 B/op      26089 allocs/op
BenchmarkGetIP-16        6268292               194.8 ns/op            64 B/op          2 allocs/op
```

## Updates

A GitHub Actions workflow is run weekly to update the IP range data if changed by the supported cloud providers. A new version is created and tagged if changes are detected. Use dependabot or renovate to automate updates to the latest version.