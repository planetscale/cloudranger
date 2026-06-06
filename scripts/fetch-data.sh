#!/usr/bin/env bash

set -eou pipefail

curl -Lo 'data/aws-ip-ranges.json' 'https://ip-ranges.amazonaws.com/ip-ranges.json'
curl -Lo 'data/gcp-cloud.json' 'https://www.gstatic.com/ipranges/cloud.json'
curl -Lo 'data/cloudflare-ips.json' 'https://api.cloudflare.com/client/v4/ips?networks=jdcloud'