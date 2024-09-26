package adapter

import (
	"context"

	"github.com/alexandremahdhaoui/dyncidr-controller/internal/types"
	"github.com/alexandremahdhaoui/dyncidr-controller/pkg/apis/v1alpha1"
)

type ResourceMutator struct{}

func (rm *ResourceMutator) Mutate(
	ctx context.Context,
	dn *v1alpha1.DynamicNetwork,
	allocatedSubnets []types.AllocatedSubnet,
) error {
	panic("not implemented")
}
