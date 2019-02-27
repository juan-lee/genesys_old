// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package bootstrap

import (
	"context"
	"github.com/juan-lee/genesys/bootstrap/cluster"
)

// Injectors from wire.go:

func initializeTestCluster(ctx context.Context) (*cluster.Bootstrapper, error) {
	bootstrapper := cluster.ProvideBootstrapper()
	return bootstrapper, nil
}
