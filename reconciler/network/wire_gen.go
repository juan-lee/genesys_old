// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package network

import (
	"context"
)

// Injectors from inject_fakes.go:

func InjectFakeReconciler(ctx context.Context) (Reconciler, error) {
	vNetReconciler := ProvideFakeVirtualNetwork()
	networkReconciler := ProvideReconciler(vNetReconciler)
	return networkReconciler, nil
}
