package providers

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	ua "github.com/wux1an/fake-useragent"
)

const (
	AWSProviderName          = "aws"
	AzureProviderName        = "azure"
	GCPProviderName          = "gcp"
	OracleProviderName       = "oracle"
	LinodeProviderName       = "linode"
	CloudflareProviderName   = "cloudflare"
	DigitalOceanProviderName = "digitalocean"
	FastlyProviderName       = "fastly"
)

type IProvider interface {
	Download() error
	Generate(f *os.File) error
}

type Provider struct {
	Output string
	URL    string
}

func (p *Provider) Download() error {
	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		return err
	}

	u, err := url.Parse(p.URL)
	if err != nil {
		return err
	}

	req.Header.Set("Host", u.Host)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", ua.Random())
	req.Header.Set("Referer", p.URL)

	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Transport: &http.Transport{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(p.Output, data, 0644)
}
