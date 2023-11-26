# cloudranger

`cloudranger` is a Go library designed to identify cloud provider information from IP addresses, with current support for AWS and GCP.

It functions without any external runtime dependencies, as IP range data is stored internally. Meant for high throughput, low-latency environments, `cloudranger` also focuses on rapid startup, loading in under 4ms. You can verify this on your system by running `make bench` and checking the `BenchmarkNew` results.

New releases are automatically created in response to updates in cloud providers' IP range information. This process, facilitated through GitHub Actions, is executed weekly to ensure the library remains up-to-date.

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

## Updates

A GitHub Actions workflow is run weekly to update the IP range data if changed by the supported cloud providers. A new version is created and tagged if changes are detected. Use dependabot or renovate to automate updates to the latest version.