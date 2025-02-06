package providers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type CloudflareProvider struct {
	*Provider
}

func (p *CloudflareProvider) Generate(out *os.File) error {
	if err := p.Download(); err != nil {
		return err
	}

	d, err := os.Open(p.Output)
	if err != nil {
		return err
	}
	defer d.Close()

	scanner := bufio.NewScanner(d)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "\t// provider: %s, cidr: %s\n", CloudflareProviderName, line)
		fmt.Fprintf(out,
			"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
			ipnet.IP[0], ipnet.IP[1], ipnet.IP[2], ipnet.IP[3],
			ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3],
			CloudflareProviderName,
			"",
		)
	}

	return nil
}

func NewCloudflareProvider(output string) *CloudflareProvider {
	return &CloudflareProvider{
		Provider: &Provider{
			Output: output,
			URL:    "https://www.cloudflare.com/ips-v4/",
		},
	}
}
