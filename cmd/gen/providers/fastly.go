package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type FastlyIPRanges struct {
	Addresses     []string `json:"addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
}

type FastlyProvider struct {
	*Provider
}

func (p *FastlyProvider) Generate(out *os.File) error {
	if err := p.Download(); err != nil {
		return err
	}

	var data FastlyIPRanges

	d, err := os.ReadFile(p.Output)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(d, &data); err != nil {
		return err
	}

	for _, ipPrefix := range append(data.Addresses, data.IPv6Addresses...) {
		_, ipnet, err := net.ParseCIDR(ipPrefix)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "\t// provider: %s, cidr: %s\n", FastlyProviderName, ipPrefix)

		if len(ipnet.IP) == net.IPv4len {
			fmt.Fprintf(out,
				"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
				ipnet.IP[0], ipnet.IP[1], ipnet.IP[2], ipnet.IP[3],
				ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3],
				FastlyProviderName,
				"",
			)
		}

		if len(ipnet.IP) == net.IPv6len {
			i := ipnet.IP
			m := ipnet.Mask
			fmt.Fprintf(out,
				"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
				i[0], i[1], i[2], i[3], i[4], i[5], i[6], i[7], i[8], i[9], i[10], i[11], i[12], i[13], i[14], i[15],
				m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15],
				FastlyProviderName,
				"",
			)
		}
	}

	return nil
}

func NewFastlyProvider(output string) *FastlyProvider {
	return &FastlyProvider{
		Provider: &Provider{
			Output: output,
			URL:    "https://api.fastly.com/public-ip-list",
		},
	}
}
