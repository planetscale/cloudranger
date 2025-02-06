package providers

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
)

type DigitalOceanIPInfo struct {
	IPPrefix   string `csv:"ip_prefix"`
	Alpha2Code string `csv:"alpha2code"`
	Region     string `csv:"region"`
	City       string `csv:"city"`
	PostalCode string `csv:"postal_code"`
}

type DigitalOceanProvider struct {
	*Provider
}

func getField[T string](record []T, index int) T {
	if len(record) > index {
		return record[index]
	}
	return ""
}

func (p *DigitalOceanProvider) Generate(out *os.File) error {
	if err := p.Download(); err != nil {
		return err
	}

	d, err := os.Open(p.Output)
	if err != nil {
		return err
	}
	defer d.Close()

	reader := csv.NewReader(d)
	reader.Comment = '#'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		ipInfo := DigitalOceanIPInfo{
			IPPrefix:   getField(record, 0),
			Alpha2Code: getField(record, 1),
			Region:     getField(record, 2),
			City:       getField(record, 3),
			PostalCode: getField(record, 4),
		}

		_, ipnet, err := net.ParseCIDR(ipInfo.IPPrefix)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "\t// provider: %s, cidr: %s, region: %s\n", DigitalOceanProviderName, ipInfo.IPPrefix, ipInfo.Region)

		if len(ipnet.IP) == net.IPv4len {
			fmt.Fprintf(out,
				"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
				ipnet.IP[0], ipnet.IP[1], ipnet.IP[2], ipnet.IP[3],
				ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3],
				DigitalOceanProviderName,
				ipInfo.Region,
			)
		}

		if len(ipnet.IP) == net.IPv6len {
			i := ipnet.IP
			m := ipnet.Mask
			fmt.Fprintf(out,
				"\t{&net.IPNet{IP: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}, Mask: []byte{%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d}}, IPInfo{cloud: \"%s\", region: \"%s\"}},\n",
				i[0], i[1], i[2], i[3], i[4], i[5], i[6], i[7], i[8], i[9], i[10], i[11], i[12], i[13], i[14], i[15],
				m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15],
				DigitalOceanProviderName,
				ipInfo.Region,
			)
		}
	}

	return nil
}

func NewDigitalOceanProvider(output string) *DigitalOceanProvider {
	return &DigitalOceanProvider{
		Provider: &Provider{
			Output: output,
			URL:    "https://www.digitalocean.com/geo/google.csv",
		},
	}
}
