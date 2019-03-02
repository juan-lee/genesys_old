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

package azure

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/google/wire"
	"github.com/juan-lee/genesys/pkg/reconcile/cluster"
	"github.com/juan-lee/genesys/pkg/reconcile/network"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func InjectCluster(c Configuration) (reconcile.Reconciler, error) {
	panic(wire.Build(
		provideAuthorizer,
		ProvideNetwork,
		network.ProvideReconciler,
		cluster.ProvideReconciler,
	))
}

func provideAuthorizer(c Configuration) (autorest.Authorizer, error) {
	env, err := azure.EnvironmentFromName(c.Cloud)
	if err != nil {
		return nil, err
	}
	oauthConfig, err := adal.NewOAuthConfig(
		env.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}

	token, err := adal.NewServicePrincipalToken(*oauthConfig, c.ClientID, c.ClientSecret, env.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}
	return autorest.NewBearerAuthorizer(token), nil
}
