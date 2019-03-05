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
	"reflect"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

var _ provider.FlatNetwork = &FlatNetwork{}

type FlatNetwork struct {
	log        logr.Logger
	config     v1alpha1.Cloud
	names      *names
	vnetClient network.VirtualNetworksClient
	sgClient   network.SecurityGroupsClient
	rtClient   network.RouteTablesClient
}

func ProvideFlatNetwork(log logr.Logger, a autorest.Authorizer, c v1alpha1.Cloud, n *names) (*FlatNetwork, error) {
	vnetClient, err := newVNETClient(c.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	sgClient, err := newNSGClient(c.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	rtClient, err := newRouteTableClient(c.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	return &FlatNetwork{
		log:        log,
		config:     c,
		names:      n,
		vnetClient: vnetClient,
		sgClient:   sgClient,
		rtClient:   rtClient,
	}, nil
}

func (r *FlatNetwork) Ensure(ctx context.Context, net v1alpha1.Network) (*provider.Status, error) {
	vnet, err := r.vnetClient.Get(ctx, r.config.ResourceGroup, r.names.VirtualNetwork(), "")
	if err != nil && IsNotFound(err) {
		if err := r.create(ctx, net); err != nil {
			return &provider.Status{}, err
		}
		return &provider.Status{State: provider.Provisioning}, nil
	} else if err != nil {
		return &provider.Status{}, err
	}

	if vnet.ProvisioningState != nil && *vnet.ProvisioningState == "Provisioning" {
		return &provider.Status{State: provider.Provisioning}, nil
	}

	if !reflect.DeepEqual(net, convert(&vnet)) {
		vnet.Location = &r.config.Location
		// TODO: probably dangerous without nil checks
		(*vnet.AddressSpace.AddressPrefixes)[0] = net.CIDR
		(*vnet.AddressSpace.AddressPrefixes)[0] = net.CIDR
		(*vnet.Subnets)[0].Name = to.StringPtr(r.names.Subnet())
		(*vnet.Subnets)[0].AddressPrefix = &net.SubnetCIDR
		_, err := r.vnetClient.CreateOrUpdate(ctx, r.config.ResourceGroup, r.names.VirtualNetwork(), vnet)
		if err != nil {
			return &provider.Status{}, err
		}
		return &provider.Status{State: provider.Provisioning}, nil
	}

	return &provider.Status{State: provider.Succeeded}, nil
}

func (r *FlatNetwork) EnsureDeleted(ctx context.Context, net v1alpha1.Network) (*provider.Status, error) {
	_, err := r.vnetClient.Delete(ctx, r.config.ResourceGroup, r.names.VirtualNetwork())
	if err != nil {
		return &provider.Status{}, err
	}
	return &provider.Status{}, nil
}

func (r *FlatNetwork) create(ctx context.Context, net v1alpha1.Network) error {
	sgf, err := r.sgClient.CreateOrUpdate(ctx, r.config.ResourceGroup, r.names.NetworkSecurityGroup(), network.SecurityGroup{
		Location: &r.config.Location,
	})
	if err != nil {
		return err
	}

	err = sgf.WaitForCompletionRef(ctx, r.sgClient.Client)
	if err != nil {
		return err
	}

	sg, err := sgf.Result(r.sgClient)
	if err != nil {
		return err
	}

	rtf, err := r.rtClient.CreateOrUpdate(ctx, r.config.ResourceGroup, r.names.RouteTable(), network.RouteTable{})
	if err != nil {
		return err
	}

	err = rtf.WaitForCompletionRef(ctx, r.rtClient.Client)
	if err != nil {
		return err
	}

	rt, err := rtf.Result(r.rtClient)
	if err != nil {
		return err
	}

	f, err := r.vnetClient.CreateOrUpdate(ctx, r.config.ResourceGroup, r.names.VirtualNetwork(),
		network.VirtualNetwork{
			Location: &r.config.Location,
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{net.CIDR},
				},
				Subnets: &[]network.Subnet{
					{
						Name: to.StringPtr(r.names.Subnet()),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix:        to.StringPtr(net.SubnetCIDR),
							NetworkSecurityGroup: &sg,
							RouteTable:           &rt,
						},
					},
				},
			},
		})
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, r.vnetClient.Client)
	if err != nil {
		return err
	}
	return err
}
