#!/usr/bin/env bash

set -eou pipefail

curl -Lo 'data/aws-ip-ranges.json' 'https://ip-ranges.amazonaws.com/ip-ranges.json'
curl -Lo 'data/gcp-cloud.json' 'https://www.gstatic.com/ipranges/cloud.json'