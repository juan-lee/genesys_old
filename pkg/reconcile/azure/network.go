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
	"errors"
	"fmt"

	aznet "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/juan-lee/genesys/pkg/reconcile/network"
)

const (
	defaultVirtualNetworkCIDR = "10.0.0.0/8"
	defaultSubnetCIDR         = "10.240.0.0/12"
	defaultSubnetName         = "k8s-subnet"
)

var _ network.VNETProvider = &netProvider{}

type netProvider struct {
	config Configuration
	client aznet.VirtualNetworksClient
}

func ProvideNetwork(c Configuration, a autorest.Authorizer) (network.VNETProvider, error) {
	client, err := newVNETClient(c, a)
	if err != nil {
		return nil, err
	}
	return &netProvider{
		config: c,
		client: client,
	}, nil
}

func (r *netProvider) State(ctx context.Context, desired network.VNETOptions) error {
	vnet, err := r.client.Get(ctx, r.config.ResourceGroup, r.vnetName(), "")
	if err != nil {
		return err
	}

	if !r.reachedDesiredState(&desired, &vnet) {
		return errors.New("OutOfSync")
	}

	return nil
}

func (r *netProvider) Update(ctx context.Context, desired network.VNETOptions) error {
	err := r.validate(&desired)
	if err != nil {
		return err
	}

	_, err = r.client.CreateOrUpdate(ctx, r.config.ResourceGroup, r.vnetName(),
		aznet.VirtualNetwork{
			Location: to.StringPtr(r.config.Location),
			VirtualNetworkPropertiesFormat: &aznet.VirtualNetworkPropertiesFormat{
				AddressSpace: &aznet.AddressSpace{
					AddressPrefixes: &[]string{desired.CIDR},
				},
				Subnets: &[]aznet.Subnet{
					{
						Name: to.StringPtr(defaultSubnetName),
						SubnetPropertiesFormat: &aznet.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr(desired.SubnetCIDR),
						},
					},
				},
			},
		})
	return err
}

func (r *netProvider) Delete(ctx context.Context) error {
	return nil
}

func (r *netProvider) vnetName() string {
	return fmt.Sprintf("%s-vnet", r.config.ResourceGroup)
}

func (r *netProvider) validate(desired *network.VNETOptions) error {
	if desired == nil {
		return errors.New("desired cannot be nil")
	}

	// TODO: validate cidr
	if desired.CIDR == "" {
		return errors.New("CIDR cannot be empty")
	}

	// TODO: validate cidr
	if desired.SubnetCIDR == "" {
		return errors.New("SubnetCIDR cannot be empty")
	}

	return nil
}

func (r *netProvider) reachedDesiredState(desired *network.VNETOptions, current *aznet.VirtualNetwork) bool {
	if desired == nil || current == nil {
		return false
	}

	if current.Location == nil || *current.Location != r.config.Location {
		return false
	}

	if current.VirtualNetworkPropertiesFormat == nil ||
		current.VirtualNetworkPropertiesFormat.AddressSpace == nil ||
		current.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes == nil ||
		len(*current.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes) != 1 ||
		current.VirtualNetworkPropertiesFormat.Subnets == nil ||
		len(*current.VirtualNetworkPropertiesFormat.Subnets) != 1 ||
		(*current.VirtualNetworkPropertiesFormat.Subnets)[0].Name == nil ||
		(*current.VirtualNetworkPropertiesFormat.Subnets)[0].AddressPrefix == nil {
		return false
	}

	if (*current.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes)[0] != desired.CIDR {
		return false
	}

	if *(*current.VirtualNetworkPropertiesFormat.Subnets)[0].Name != defaultSubnetName ||
		*(*current.VirtualNetworkPropertiesFormat.Subnets)[0].AddressPrefix != desired.SubnetCIDR {
		return false
	}

	return false
}
