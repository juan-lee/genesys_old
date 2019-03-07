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
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

func (p *Provider) GetControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) (exists bool, err error) {
	_, err = p.client.pip.Get(ctx, p.config.ResourceGroup, p.names.ControlPlaneEndpoint(), "")
	if err != nil && IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (p *Provider) EnsureControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) error {
	f, err := p.client.pip.CreateOrUpdate(ctx, p.config.ResourceGroup, p.names.ControlPlaneEndpoint(), network.PublicIPAddress{
		Sku: &network.PublicIPAddressSku{Name: network.PublicIPAddressSkuNameBasic},
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   network.IPv4,
			PublicIPAllocationMethod: network.Static,
			DNSSettings: &network.PublicIPAddressDNSSettings{
				DomainNameLabel: to.StringPtr(p.names.ControlPlaneEndpoint()),
				Fqdn:            to.StringPtr(cp.Fqdn),
			},
		},
	})
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, p.client.pip.Client)
	if err != nil {
		return err
	}

	_, err = f.Result(p.client.pip)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) UpdateControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) error {
	pip, err := p.client.pip.Get(ctx, p.config.ResourceGroup, p.names.ControlPlaneEndpoint(), "")
	if err != nil && IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	// TODO: improve this conversion
	pip.Location = &p.config.ResourceGroup
	pip.Sku = &network.PublicIPAddressSku{Name: network.PublicIPAddressSkuNameBasic}
	pip.PublicIPAddressVersion = network.IPv4
	pip.PublicIPAllocationMethod = network.Static
	pip.DNSSettings.DomainNameLabel = to.StringPtr(p.names.ControlPlaneEndpoint())
	pip.DNSSettings.Fqdn = to.StringPtr(p.names.ControlPlaneEndpoint())
	f, err := p.client.pip.CreateOrUpdate(ctx, p.config.ResourceGroup, p.names.ControlPlaneEndpoint(), pip)
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, p.client.pip.Client)
	if err != nil {
		return err
	}

	_, err = f.Result(p.client.pip)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureControlPlaneEndpointDeleted(ctx context.Context, cp *v1alpha1.ControlPlane) error {
	f, err := p.client.pip.Delete(ctx, p.config.ResourceGroup, p.names.ControlPlaneEndpoint())
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, p.client.pip.Client)
	if err != nil {
		return err
	}

	_, err = f.Result(p.client.pip)
	if err != nil {
		return err
	}

	return nil
}
