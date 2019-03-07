// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package azure

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/go-logr/logr"
	"github.com/google/wire"
	"github.com/juan-lee/genesys/pkg/actuator/cluster"
	"github.com/juan-lee/genesys/pkg/actuator/controlplane"
	"github.com/juan-lee/genesys/pkg/actuator/network"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	"github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// Injectors from inject_azure.go:

func InjectCluster(log logr.Logger, cloud *v1alpha1.Cloud) (*cluster.SelfManaged, error) {
	configuration, err := provideConfiguration()
	if err != nil {
		return nil, err
	}
	authorizer, err := provideAuthorizer(log, configuration)
	if err != nil {
		return nil, err
	}
	azureNames := provideNames(cloud)
	azureVirtualNetworkFactory, err := provideVirtualNetworkFactory(log, authorizer, cloud, azureNames)
	if err != nil {
		return nil, err
	}
	flat, err := network.ProvideFlatNetwork(log, azureVirtualNetworkFactory)
	if err != nil {
		return nil, err
	}
	azureControlPlaneEndpointFactory, err := provideControlPlaneEndpointFactory(log, authorizer, cloud, azureNames)
	if err != nil {
		return nil, err
	}
	singleInstance, err := controlplane.ProvideSingleInstance(log, azureControlPlaneEndpointFactory)
	if err != nil {
		return nil, err
	}
	selfManaged, err := cluster.ProvideSelfManaged(log, flat, singleInstance)
	if err != nil {
		return nil, err
	}
	return selfManaged, nil
}

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

var cpSet = wire.NewSet(
	provideControlPlaneEndpointFactory, wire.Bind(new(provider.ControlPlaneEndpointFactory), new(controlPlaneEndpointFactory)),
)

var netSet = wire.NewSet(
	provideVirtualNetworkFactory, wire.Bind(new(provider.VirtualNetworkFactory), new(virtualNetworkFactory)),
)

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
