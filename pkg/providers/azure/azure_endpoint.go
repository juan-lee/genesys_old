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

package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

var _ provider.Provider = &controlPlaneEndpoint{}
var _ provider.ControlPlaneEndpointFactory = &controlPlaneEndpointFactory{}

type controlPlaneEndpointFactory struct {
	log    logr.Logger
	config *v1alpha1.Cloud
	names  *names
	auth   autorest.Authorizer
}

type controlPlaneEndpoint struct {
	log     logr.Logger
	config  *v1alpha1.Cloud
	names   *names
	pip     *network.PublicIPAddress
	desired *v1alpha1.ControlPlane
	client  network.PublicIPAddressesClient
}

func provideControlPlaneEndpointFactory(log logr.Logger, a autorest.Authorizer, c *v1alpha1.Cloud, n *names) (*controlPlaneEndpointFactory, error) {
	return &controlPlaneEndpointFactory{
		log:    log,
		config: c,
		names:  n,
		auth:   a,
	}, nil
}

func (f *controlPlaneEndpointFactory) Get(ctx context.Context, cp *v1alpha1.ControlPlane) (provider.Reconciler, error) {
	pip, err := newControlPlaneEndpoint(ctx, f, cp)
	if err != nil {
		return nil, err
	}
	return &Reconciler{log: f.log, Provider: pip}, nil
}

func newControlPlaneEndpoint(ctx context.Context, f *controlPlaneEndpointFactory, cp *v1alpha1.ControlPlane) (*controlPlaneEndpoint, error) {
	client, err := newPublicIPClient(f.config.SubscriptionID, f.auth)
	if err != nil {
		return nil, err
	}
	var current *network.PublicIPAddress
	if pip, err := client.Get(ctx, f.config.ResourceGroup, f.names.ControlPlaneEndpoint(), ""); err == nil {
		current = &pip
	} else if err != nil && !IsNotFound(err) {
		return nil, err
	}
	return &controlPlaneEndpoint{
		log:     f.log,
		config:  f.config,
		names:   f.names,
		pip:     current,
		desired: cp,
		client:  client,
	}, nil
}

func (c *controlPlaneEndpoint) Exists() bool {
	if c.pip == nil {
		return false
	}
	return true
}

func (c *controlPlaneEndpoint) Status() provider.Status {
	if !c.Exists() {
		return provider.NeedsUpdate
	}

	if c.pip.ProvisioningState == nil {
		return provider.Unknown
	}

	switch *c.pip.ProvisioningState {
	case "Succeeded":
		return provider.Succeeded
	case "Provisioning":
		return provider.Provisioning
	case "Deleting":
		return provider.Provisioning
	}

	return provider.Unknown
}

func (c *controlPlaneEndpoint) Update(ctx context.Context) error {
	if c.Exists() {
		c.pip.Location = &c.config.ResourceGroup
		c.pip.Sku = &network.PublicIPAddressSku{Name: network.PublicIPAddressSkuNameBasic}
		c.pip.PublicIPAddressVersion = network.IPv4
		c.pip.PublicIPAllocationMethod = network.Static
		c.pip.DNSSettings.DomainNameLabel = to.StringPtr(c.names.ControlPlaneEndpoint())
		c.pip.DNSSettings.Fqdn = to.StringPtr(c.names.ControlPlaneEndpoint())
		_, err := c.client.CreateOrUpdate(ctx, c.config.ResourceGroup, c.names.ControlPlaneEndpoint(), *c.pip)
		if err != nil {
			return err
		}
		return nil
	}

	f, err := c.client.CreateOrUpdate(ctx, c.config.ResourceGroup, c.names.ControlPlaneEndpoint(), network.PublicIPAddress{
		Sku: &network.PublicIPAddressSku{Name: network.PublicIPAddressSkuNameBasic},
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   network.IPv4,
			PublicIPAllocationMethod: network.Static,
			DNSSettings: &network.PublicIPAddressDNSSettings{
				DomainNameLabel: to.StringPtr(c.names.ControlPlaneEndpoint()),
				Fqdn:            to.StringPtr(c.desired.Fqdn),
			},
		},
	})
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, c.client.Client)
	if err != nil {
		return err
	}

	return nil
}

func (c *controlPlaneEndpoint) Delete(ctx context.Context) error {
	panic("not implemented")
}
