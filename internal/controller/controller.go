package controller

import (
	"context"

	"github.com/alexandremahdhaoui/dyncidr-controller/internal/adapter"
	"github.com/alexandremahdhaoui/dyncidr-controller/pkg/apis/v1alpha1"
)

type Reconciler struct {
	NetworkFetcher  adapter.NetworkFetcher
	ResourceMutator adapter.ResourceMutator
	SubnetAllocator adapter.SubnetAllocator
}

func (r *Reconciler) Reconcile(ctx context.Context, dn *v1alpha1.DynamicNetwork) {
	ipNet, err := r.NetworkFetcher.Fetch(ctx, dn)
	if err != nil {
		return nil, err
	}

	allocatedSubnets, err := r.SubnetAllocator.Allocate(dn, ipNet)
	if err != nil {
		return nil, err
	}

	r.ResourceMutator.Mutate(ctx, dn, allocatedSubnets)
}
