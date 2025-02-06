package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type AzureIPData struct {
	ChangeNumber int            `json:"changeNumber"`
	Cloud        string         `json:"cloud"`
	Values       []AzureIPValue `json:"values"`
}

type AzureIPValue struct {
	Name       string                `json:"name"`
	ID         string                `json:"id"`
	Properties AzureRegionProperties `json:"properties"`
}

type AzureRegionProperties struct {
	ChangeNumber    int      `json:"changeNumber"`
	Region          string   `json:"region"`
	RegionID        int      `json:"regionId"`
	Platform        string   `json:"platform"`
	SystemService   string   `json:"systemService"`
	AddressPrefixes []string `json:"addressPrefixes"`
}

type AzureProvider struct {
	*Provider
}

func (p *AzureProvider) Generate(out *os.File) error {
	if err := p.Download(); err != nil {
		return err
	}

	d, err := os.ReadFile(p.Output)
	if err != nil {
		return err
	}

	var data AzureIPData

	if err := json.Unmarshal(d, &data); err != nil {
		return err
	}

	for _, value := range data.Values {
		for _, prefix := range value.Properties.AddressPrefixes {
			_, ipnet, err := net.ParseCIDR(prefix)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "\t// provider: %s, cidr: %s, region: %s\n", AzureProviderName, prefix, value.Properties.Region)

			if len(ipnet.IP) == net.IPv4len {
				fmt.Fprintf(out,
					"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
					ipnet.IP[0], ipnet.IP[1], ipnet.IP[2], ipnet.IP[3],
					ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3],
					AzureProviderName,
					value.Properties.Region,
				)
			}

			if len(ipnet.IP) == net.IPv6len {
				i := ipnet.IP
				m := ipnet.Mask
				fmt.Fprintf(out,
					"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
					i[0], i[1], i[2], i[3], i[4], i[5], i[6], i[7], i[8], i[9], i[10], i[11], i[12], i[13], i[14], i[15],
					m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15],
					AzureProviderName,
					value.Properties.Region,
				)
			}
		}
	}

	return nil
}

func NewAzureProvider(output string) *AzureProvider {
	return &AzureProvider{
		Provider: &Provider{
			Output: output,
			// URL:    "https://www.microsoft.com/en-us/download/details.aspx?id=56519",
			URL: "https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_20250203.json",
		},
	}
}
