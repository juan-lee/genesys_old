// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package azure

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// Injectors from inject_azure.go:

func injectProvider(cloud *v1alpha1.Cloud) (*Provider, error) {
	logger, err := provideLogger()
	if err != nil {
		return nil, err
	}
	azureNames := provideNames(cloud)
	configuration, err := provideConfiguration()
	if err != nil {
		return nil, err
	}
	authorizer, err := provideAuthorizer(logger, configuration)
	if err != nil {
		return nil, err
	}
	azureClient, err := provideClient(logger, cloud, authorizer)
	if err != nil {
		return nil, err
	}
	provider, err := setupProvider(logger, cloud, azureNames, azureClient)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// inject_azure.go:

func provideLogger() (logr.Logger, error) {
	return log.Log.WithName("azure.provider"), nil
}

func provideClient(log2 logr.Logger, cloud *v1alpha1.Cloud, a autorest.Authorizer) (*client, error) {
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

func setupProvider(log2 logr.Logger, cloud *v1alpha1.Cloud, names2 *names, client2 *client) (*Provider, error) {
	return &Provider{
		log:    log2,
		config: cloud,
		names:  names2,
		client: client2,
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

func provideAuthorizer(log2 logr.Logger, config *Configuration) (autorest.Authorizer, error) {
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
