# cloudranger

`cloudranger` is a Go library designed to identify cloud provider information from IP addresses, with current support for AWS and GCP.

It functions without any external runtime dependencies, as IP range data is stored internally. Meant for high throughput, low-latency environments, `cloudranger` also focuses on rapid startup, loading in under 4ms. You can verify this on your system by running `make bench` and checking the `BenchmarkNew` results.

New releases are automatically created in response to updates in cloud providers' IP range information. This process, facilitated through GitHub Actions, is executed weekly to ensure the library remains up-to-date.

The inspiration for `cloudranger` came from a similar library found at https://github.com/kubernetes/registry.k8s.io, used by the Kubernetes OCI registry for redirecting requests to the appropriate cloud provider. We developed `cloudranger` to provide a standalone library adaptable for various projects, offering greater control for our specific use cases and minimizing the impact of upstream changes. Unlike the original project, which uses its own trie implementation, `cloudranger` depends on github.com/infobloxopen/go-trees. While both implementations have not been directly benchmarked against each other, their performance is expected to be comparable.

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

## IP Range Database Updates

IP range data is sourced from:

- AWS: https://ip-ranges.amazonaws.com/ip-ranges.json
- GCP: https://www.gstatic.com/ipranges/cloud.json

A GitHub Actions workflow is run weekly to update the IP range data if changed by the supported cloud providers. A new version is created and tagged if changes are detected. Use dependabot or renovate to automate updates to the latest version.