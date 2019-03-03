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
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/go-logr/logr"
	"github.com/google/wire"
	k8sv1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"github.com/juan-lee/genesys/pkg/reconcile/cluster"
	"github.com/juan-lee/genesys/pkg/reconcile/network"
)

func InjectCluster(log logr.Logger, c k8sv1alpha1.Cloud) (cluster.Reconciler, error) {
	panic(wire.Build(
		provideConfiguration,
		provideAuthorizer,
		ProvideNetwork,
		network.ProvideReconciler,
		cluster.ProvideReconciler,
	))
}

func provideConfiguration() (Configuration, error) {
	return Configuration{
		Cloud:        "AzurePublicCloud",
		ClientID:     os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
		TenantID:     os.Getenv("AZURE_TENANT_ID"),
		UserAgent:    "genesys",
	}, nil
}

func provideAuthorizer(log logr.Logger, c Configuration) (autorest.Authorizer, error) {
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
