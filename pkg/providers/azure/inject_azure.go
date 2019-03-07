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
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func injectProvider(cloud *v1alpha1.Cloud) (*Provider, error) {
	panic(wire.Build(
		provideLogger,
		provideNames,
		provideConfiguration,
		provideAuthorizer,
		provideClient,
		setupProvider,
	))
}

func provideLogger() (logr.Logger, error) {
	return logf.Log.WithName("azure.provider"), nil
}

func provideClient(log logr.Logger, cloud *v1alpha1.Cloud, a autorest.Authorizer) (*client, error) {
	vnet, err := newVNETClient(cloud.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	nsg, err := newNSGClient(cloud.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	rt, err := newRouteTableClient(cloud.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	return &client{
		vnet: vnet,
		nsg:  nsg,
		rt:   rt,
	}, nil
}

func setupProvider(log logr.Logger, cloud *v1alpha1.Cloud, names *names, client *client) (*Provider, error) {
	return &Provider{
		log:    log,
		config: cloud,
		names:  names,
		client: client,
	}, nil
}

func provideNames(cloud *v1alpha1.Cloud) *names {
	return &names{prefix: cloud.ResourceGroup}
}

func provideConfiguration() (*Configuration, error) {
	return &Configuration{
		Cloud:        "AzurePublicCloud",
		ClientID:     os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
		TenantID:     os.Getenv("AZURE_TENANT_ID"),
	}, nil
}

func provideAuthorizer(log logr.Logger, config *Configuration) (autorest.Authorizer, error) {
	env, err := azure.EnvironmentFromName(config.Cloud)
	if err != nil {
		return nil, err
	}
	oauthConfig, err := adal.NewOAuthConfig(
		env.ActiveDirectoryEndpoint, config.TenantID)
	if err != nil {
		return nil, err
	}

	token, err := adal.NewServicePrincipalToken(*oauthConfig, config.ClientID, config.ClientSecret, env.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}
	return autorest.NewBearerAuthorizer(token), nil
}
