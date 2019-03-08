// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//+build wireinject

package bootstrap

import (
	"github.com/go-logr/logr"
	"github.com/google/wire"
	"github.com/juan-lee/genesys/pkg/actuator/controlplane"
	"github.com/juan-lee/genesys/pkg/actuator/network"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"github.com/juan-lee/genesys/pkg/providers/azure"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func injectBootstrap(cloud *v1alpha1.Cloud) (*Bootstrap, error) {
	panic(wire.Build(
		provideLogger,
		setupAzureProvider,
		network.ProvideFlat,
		controlplane.ProvideSingleInstance,
		provideBootstrap,
	))
}

func provideLogger() (logr.Logger, error) {
	return logf.Log.WithName("actuator.bootstrap"), nil
}

func setupAzureProvider(log logr.Logger, cloud *v1alpha1.Cloud) (provider.Interface, error) {
	return azure.NewProvider(cloud)
}

func provideBootstrap(log logr.Logger, net *network.Flat, cp *controlplane.SingleInstance) (*Bootstrap, error) {
	return &Bootstrap{log: log, net: net, cp: cp}, nil
}
