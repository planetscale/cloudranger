package cloudranger

import (
	"net"

	"github.com/infobloxopen/go-trees/iptree"
)

type CloudRanger struct {
	tree *iptree.Tree
}

// IPInfo contains information about an IP address.
type IPInfo struct {
	cloud  string
	region string
}

// New returns a new cloudRanger.
func New() *CloudRanger {
	tree := iptree.NewTree()
	for _, r := range cloudRanges {
		tree.InplaceInsertNet(r.net, r.info)
	}
	return &CloudRanger{
		tree: tree,
	}
}

// GetIP returns the IPInfo for the given IP address. If the IP address is not
// found in any of the known cloud providers the second return value is false.
func (cr *CloudRanger) GetIP(ip string) (IPInfo, bool) {
	addr := net.ParseIP(ip)
	if addr == nil {
		return IPInfo{}, false
	}

	n, found := cr.tree.GetByIP(addr)
	if !found {
		return IPInfo{}, false
	}
	return n.(IPInfo), true
}

// Cloud returns the cloud provider for the IP address.
func (i IPInfo) Cloud() string {
	return i.cloud
}

// Region returns the cloud provider's region for the IP address.
func (i IPInfo) Region() string {
	return i.region
}
