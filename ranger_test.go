package cloudranger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	cr := New()

	tests := []struct {
		name           string
		ip             string
		expectedCloud  string
		expectedRegion string
		found          bool
	}{
		{
			name:           "valid IPv4 address in Amazon Web Services",
			ip:             "3.5.140.101",
			expectedCloud:  "AWS",
			expectedRegion: "ap-northeast-2",
			found:          true,
		},
		{
			name:           "valid IPv6 address in Amazon Web Services",
			ip:             "2a05:d077:6081::1",
			expectedCloud:  "AWS",
			expectedRegion: "eu-north-1",
			found:          true,
		},
		{
			name:           "valid IPv4 address in Google Cloud Platform",
			ip:             "34.35.1.2",
			expectedCloud:  "GCP",
			expectedRegion: "africa-south1",
			found:          true,
		},
		{
			name:           "valid IPv6 address in Google Cloud Platform",
			ip:             "2600:1900:4010::0000:1",
			expectedCloud:  "GCP",
			expectedRegion: "europe-west1",
			found:          true,
		},
		{
			name:  "non-cloud IP address",
			ip:    "127.0.0.1",
			found: false,
		},
		{
			name:  "not an IP address",
			ip:    "just a random string",
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipinfo, found := cr.GetIP(tt.ip)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.expectedCloud, ipinfo.Cloud())
			assert.Equal(t, tt.expectedRegion, ipinfo.Region())
		})
	}
}

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = New()
	}
}

func BenchmarkGetIP(b *testing.B) {
	b.StopTimer()
	ranger := New()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, _ = ranger.GetIP("34.35.1.3")
	}
}
