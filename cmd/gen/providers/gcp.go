package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// GCPIPData represents the structure of the JSON data from AWS.
type GCPIPData struct {
	SyncToken  string      `json:"syncToken"`
	CreateDate string      `json:"createDate"`
	Prefixes   []GCPPrefix `json:"prefixes"`
}

// GCPPrefix represents an individual IP prefix entry.
type GCPPrefix struct {
	IPPrefix   string `json:"ipv4Prefix"`
	IPv6Prefix string `json:"ipv6Prefix"`
	Scope      string `json:"scope"`   // scope is the region on GCP
	Service    string `json:"service"` // currently ignored, and it's always the same 'Google Cloud Platform' anyway
}

type GCPProvider struct {
	*Provider
}

func (p *GCPProvider) Generate(out *os.File) error {
	if err := p.Download(); err != nil {
		return err
	}

	var data GCPIPData

	d, err := os.ReadFile(p.Output)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(d, &data); err != nil {
		return err
	}

	for _, prefix := range data.Prefixes {
		var ipPrefix string
		if prefix.IPPrefix != "" {
			ipPrefix = prefix.IPPrefix
		}

		if prefix.IPv6Prefix != "" {
			ipPrefix = prefix.IPv6Prefix
		}

		if ipPrefix == "" {
			continue
		}

		_, ipnet, err := net.ParseCIDR(ipPrefix)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "\t// provider: %s, cidr: %s, region: %s\n", GCPProviderName, ipPrefix, prefix.Scope)

		if len(ipnet.IP) == net.IPv4len {
			fmt.Fprintf(out,
				"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
				ipnet.IP[0], ipnet.IP[1], ipnet.IP[2], ipnet.IP[3],
				ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3],
				GCPProviderName,
				prefix.Scope,
			)
		}

		if len(ipnet.IP) == net.IPv6len {
			i := ipnet.IP
			m := ipnet.Mask
			fmt.Fprintf(out,
				"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
				i[0], i[1], i[2], i[3], i[4], i[5], i[6], i[7], i[8], i[9], i[10], i[11], i[12], i[13], i[14], i[15],
				m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15],
				GCPProviderName,
				prefix.Scope,
			)
		}
	}

	return nil
}

func NewGCPProvider(output string) *GCPProvider {
	return &GCPProvider{
		Provider: &Provider{
			Output: output,
			URL:    "https://www.gstatic.com/ipranges/cloud.json",
		},
	}
}
