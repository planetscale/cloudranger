package providers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type CloudflareIPv6Provider struct {
	*Provider
}

func (p *CloudflareIPv6Provider) Generate(out *os.File) error {
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

		i := ipnet.IP
		m := ipnet.Mask
		fmt.Fprintf(out, "\t// provider: %s, cidr: %s\n", CloudflareProviderName, line)
		fmt.Fprintf(out,
			"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
			i[0], i[1], i[2], i[3], i[4], i[5], i[6], i[7], i[8], i[9], i[10], i[11], i[12], i[13], i[14], i[15],
			m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15],
			CloudflareProviderName,
			"",
		)
	}

	return nil
}

func NewCloudflareIPv6Provider(output string) *CloudflareIPv6Provider {
	return &CloudflareIPv6Provider{
		Provider: &Provider{
			Output: output,
			URL:    "https://www.cloudflare.com/ips-v6/",
		},
	}
}
