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

var _ provider.VirtualNetwork = &VirtualNetwork{}

type VirtualNetworkFactory struct {
	log    logr.Logger
	config *v1alpha1.Cloud
	names  *names
	auth   autorest.Authorizer
}

type VirtualNetwork struct {
	log        logr.Logger
	config     *v1alpha1.Cloud
	names      *names
	vnet       *network.VirtualNetwork
	desired    *v1alpha1.Network
	vnetClient network.VirtualNetworksClient
	sgClient   network.SecurityGroupsClient
	rtClient   network.RouteTablesClient
}

func ProvideVirtualNetworkFactory(log logr.Logger, a autorest.Authorizer, c *v1alpha1.Cloud, n *names) (*VirtualNetworkFactory, error) {
	return &VirtualNetworkFactory{
		log:    log,
		config: c,
		names:  n,
	}, nil
}

func (r *VirtualNetworkFactory) Get(ctx context.Context, net *v1alpha1.Network) (provider.Reconciler, error) {
	vnet, err := NewVirtualNetwork(ctx, r, net)
	if err != nil {
		return nil, err
	}
	return &Reconciler{Provider: vnet}, nil
}

func NewVirtualNetwork(ctx context.Context, f *VirtualNetworkFactory, desired *v1alpha1.Network) (*VirtualNetwork, error) {
	vnetClient, err := newVNETClient(f.config.SubscriptionID, f.auth)
	if err != nil {
		return nil, err
	}
	sgClient, err := newNSGClient(f.config.SubscriptionID, f.auth)
	if err != nil {
		return nil, err
	}
	rtClient, err := newRouteTableClient(f.config.SubscriptionID, f.auth)
	if err != nil {
		return nil, err
	}
	var current *network.VirtualNetwork
	if vnet, err := vnetClient.Get(ctx, f.config.ResourceGroup, f.names.VirtualNetwork(), ""); err == nil {
		current = &vnet
	} else if err != nil && !IsNotFound(err) {
		return nil, err
	}
	return &VirtualNetwork{
		log:        f.log,
		config:     f.config,
		names:      f.names,
		vnet:       current,
		vnetClient: vnetClient,
		sgClient:   sgClient,
		rtClient:   rtClient,
	}, nil
}

func (r *VirtualNetwork) Exists() bool {
	if r.vnet == nil {
		return true
	}
	return false
}

func (r *VirtualNetwork) Status() provider.Status {
	if r.Exists() {
		switch *r.vnet.ProvisioningState {
		case "Succeeded":
			return provider.Succeeded
		case "Provisioning":
			return provider.Provisioning
		case "Deleting":
			return provider.Provisioning
		}
	}
	return provider.Succeeded
}

func (r *VirtualNetwork) Update(ctx context.Context) error {
	if r.Exists() {
	}

	return r.create(ctx, r.desired)
}

func (r *VirtualNetwork) Delete(ctx context.Context) error {
	_, err := r.vnetClient.Delete(ctx, r.config.ResourceGroup, r.names.VirtualNetwork())
	if err != nil {
		return err
	}
	return nil
}

func (r *VirtualNetwork) Ensure(ctx context.Context, net *v1alpha1.Network) error {
	vnet, err := r.vnetClient.Get(ctx, r.config.ResourceGroup, r.names.VirtualNetwork(), "")
	if err != nil && IsNotFound(err) {
		if err := r.create(ctx, net); err != nil {
			return err
		}
		return r.statusProvisioning()
	} else if err != nil {
		return err
	}

	if vnet.ProvisioningState != nil && *vnet.ProvisioningState == "Provisioning" {
		return r.statusProvisioning()
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
			return err
		}
		return r.statusProvisioning()
	}

	return nil
}

func (r *VirtualNetwork) EnsureDeleted(ctx context.Context, net *v1alpha1.Network) error {
	_, err := r.vnetClient.Delete(ctx, r.config.ResourceGroup, r.names.VirtualNetwork())
	if err != nil {
		return err
	}
	return nil
}

func (r *VirtualNetwork) create(ctx context.Context, net *v1alpha1.Network) error {
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

func (r *VirtualNetwork) statusProvisioning() error {
	return provider.Pending("VirtualNetwork")
}

func convert(in *network.VirtualNetwork) v1alpha1.Network {
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
