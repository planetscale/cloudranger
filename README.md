# CloudRanger

CloudRanger is a Go library designed to identify cloud provider information from IP addresses, with current support for AWS and GCP.

It functions without any external runtime dependencies, as IP range data is stored internally. Meant for high throughput, low-latency environments, CloudRanger also focuses on rapid startup, loading in under 4ms. You can verify this on your system by running `make bench` and checking the `BenchmarkNew` results.

New releases are automatically created in response to updates in cloud providers' IP range information. This process, facilitated through GitHub Actions, is executed weekly to ensure the library remains up-to-date.