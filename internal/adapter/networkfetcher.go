package adapter

import (
	"context"
	"net"

	"github.com/alexandremahdhaoui/dyncidr-controller/pkg/apis/v1alpha1"
)

type NetworkFetcher struct{}

func (nf *NetworkFetcher) Fetch(
	ctx context.Context,
	dn *v1alpha1.DynamicNetwork,
) (*net.IPNet, error) {
	panic("implement me")
}
