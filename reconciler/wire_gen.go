// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package reconciler

import (
	"context"
	"github.com/juan-lee/genesys/reconciler/cluster"
	"github.com/juan-lee/genesys/reconciler/network"
)

// Injectors from inject_fakes.go:

func InjectFakeSelfManaged(ctx context.Context) (*cluster.Reconciler, error) {
	vNetReconciler := network.ProvideFakeVirtualNetwork()
	reconciler := network.ProvideReconciler(vNetReconciler)
	clusterReconciler := cluster.ProvideSelfManaged(reconciler)
	return clusterReconciler, nil
}
