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
	"fmt"
	"reflect"

	aznet "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/go-logr/logr"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"github.com/juan-lee/genesys/pkg/reconcile/provider"
)

const (
	defaultVirtualNetworkCIDR = "10.0.0.0/8"
	defaultSubnetCIDR         = "10.240.0.0/12"
	defaultSubnetName         = "k8s-subnet"
)

var _ provider.VirtualNetwork = &VirtualNetwork{}

type VirtualNetwork struct {
	log    logr.Logger
	config v1alpha1.Cloud
	client aznet.VirtualNetworksClient
}

func ProvideVirtualNetwork(log logr.Logger, a autorest.Authorizer, c v1alpha1.Cloud) (*VirtualNetwork, error) {
	client, err := newVNETClient(c.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	return &VirtualNetwork{
		log:    log,
		config: c,
		client: client,
	}, nil
}

func (r *VirtualNetwork) Get(ctx context.Context) (*v1alpha1.Network, error) {
	vnet, err := r.client.Get(ctx, r.config.ResourceGroup, r.vnetName(), "")
	if err != nil {
		if derr, ok := err.(autorest.DetailedError); ok && derr.StatusCode == 404 {
			return nil, provider.ErrNotFound
		}
		return nil, err
	}

	result := convert(&vnet)

	return &result, nil
}

func (r *VirtualNetwork) Update(ctx context.Context, desired v1alpha1.Network) error {
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

func (r *VirtualNetwork) Delete(ctx context.Context) error {
	_, err := r.client.Delete(ctx, r.config.ResourceGroup, r.vnetName())
	if err != nil {
		return err
	}
	return nil
}

func (r *VirtualNetwork) vnetName() string {
	return fmt.Sprintf("%s-vnet", r.config.ResourceGroup)
}

func (r *VirtualNetwork) reachedDesiredState(desired *v1alpha1.Network, current *aznet.VirtualNetwork) bool {
	if current.Location == nil || *current.Location != r.config.Location {
		r.log.Info("location is not in sync", "location", current.Location)
		return false
	}

	return reflect.DeepEqual(*desired, convert(current))
}

func convert(in *aznet.VirtualNetwork) v1alpha1.Network {
	out := v1alpha1.Network{}

	if in != nil &&
		in.VirtualNetworkPropertiesFormat != nil &&
		in.VirtualNetworkPropertiesFormat.AddressSpace != nil &&
		in.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes != nil &&
		len(*in.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes) == 1 {
		out.CIDR = (*in.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes)[0]
	}

	if in != nil &&
		in.VirtualNetworkPropertiesFormat != nil &&
		in.VirtualNetworkPropertiesFormat.Subnets != nil &&
		len(*in.VirtualNetworkPropertiesFormat.Subnets) == 1 &&
		(*in.VirtualNetworkPropertiesFormat.Subnets)[0].Name != nil &&
		(*in.VirtualNetworkPropertiesFormat.Subnets)[0].AddressPrefix != nil {
		out.SubnetCIDR = *(*in.VirtualNetworkPropertiesFormat.Subnets)[0].AddressPrefix
	}

	return out
}
