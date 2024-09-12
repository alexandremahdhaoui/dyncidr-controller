package adapter

import (
	"github.com/alexandremahdhaoui/dyncidr-controller/internal/types"
	"github.com/alexandremahdhaoui/dyncidr-controller/pkg/apis/v1alpha1"
)

type SubnetAllocator struct{}

func (sa *SubnetAllocator) Allocate(dn *v1alpha1.DynamicNetwork) ([]types.AllocatedSubnet, error) {
	// - Sort subnet by size.
	// - Allocates biggest to smallest
	// - Each iteration starts from begining of the "array" of allocatable IP addresses and try to
	//   fill up the "holes".
	panic("implement me")
}
