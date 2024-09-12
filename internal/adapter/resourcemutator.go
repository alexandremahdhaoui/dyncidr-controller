package adapter

import (
	"context"
	"go/constant"
)

type ResourceMutator struct{}

func (rm *ResourceMutator) Mutate(ctx context.Context, dn, allocatedSubnets) error {}
