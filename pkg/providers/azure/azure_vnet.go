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

var _ provider.Provider = &virtualNetwork{}
var _ provider.VirtualNetworkFactory = &virtualNetworkFactory{}

type virtualNetworkFactory struct {
	log    logr.Logger
	config *v1alpha1.Cloud
	names  *names
	auth   autorest.Authorizer
}

type virtualNetwork struct {
	log        logr.Logger
	config     *v1alpha1.Cloud
	names      *names
	vnet       *network.VirtualNetwork
	desired    *v1alpha1.Network
	vnetClient network.VirtualNetworksClient
	sgClient   network.SecurityGroupsClient
	rtClient   network.RouteTablesClient
}

func provideVirtualNetworkFactory(log logr.Logger, a autorest.Authorizer, c *v1alpha1.Cloud, n *names) (*virtualNetworkFactory, error) {
	return &virtualNetworkFactory{
		log:    log,
		config: c,
		names:  n,
		auth:   a,
	}, nil
}

func (f *virtualNetworkFactory) Get(ctx context.Context, net *v1alpha1.Network) (provider.Reconciler, error) {
	vnet, err := newVirtualNetwork(ctx, f, net)
	if err != nil {
		return nil, err
	}
	return &Reconciler{log: f.log, Provider: vnet}, nil
}

func newVirtualNetwork(ctx context.Context, f *virtualNetworkFactory, desired *v1alpha1.Network) (*virtualNetwork, error) {
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
	return &virtualNetwork{
		log:        f.log,
		config:     f.config,
		names:      f.names,
		vnet:       current,
		desired:    desired,
		vnetClient: vnetClient,
		sgClient:   sgClient,
		rtClient:   rtClient,
	}, nil
}

func (vn *virtualNetwork) Exists() bool {
	if vn.vnet == nil {
		return false
	}
	return true
}

func (vn *virtualNetwork) Status() provider.Status {
	if !vn.Exists() {
		return provider.NeedsUpdate
	}

	if vn.vnet.ProvisioningState == nil {
		return provider.Unknown
	}

	switch *vn.vnet.ProvisioningState {
	case "Succeeded":
		return provider.Succeeded
	case "Provisioning":
		return provider.Provisioning
	case "Deleting":
		return provider.Provisioning
	}

	return provider.Unknown
}

func (vn *virtualNetwork) Update(ctx context.Context) error {
	if vn.Exists() {
		if !reflect.DeepEqual(*vn.desired, convert(vn.vnet)) {
			vn.vnet.Location = &vn.config.Location
			// TODO: probably dangerous without nil checks
			(*vn.vnet.AddressSpace.AddressPrefixes)[0] = vn.desired.CIDR
			(*vn.vnet.AddressSpace.AddressPrefixes)[0] = vn.desired.CIDR
			(*vn.vnet.Subnets)[0].Name = to.StringPtr(vn.names.Subnet())
			(*vn.vnet.Subnets)[0].AddressPrefix = &vn.desired.SubnetCIDR
			_, err := vn.vnetClient.CreateOrUpdate(ctx, vn.config.ResourceGroup, vn.names.VirtualNetwork(), *vn.vnet)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return vn.create(ctx, vn.desired)
}

func (vn *virtualNetwork) Delete(ctx context.Context) error {
	_, err := vn.vnetClient.Delete(ctx, vn.config.ResourceGroup, vn.names.VirtualNetwork())
	if err != nil {
		return err
	}
	return nil
}

func (vn *virtualNetwork) create(ctx context.Context, net *v1alpha1.Network) error {
	sgf, err := vn.sgClient.CreateOrUpdate(ctx, vn.config.ResourceGroup, vn.names.NetworkSecurityGroup(), network.SecurityGroup{
		Location: &vn.config.Location,
	})
	if err != nil {
		return err
	}

	err = sgf.WaitForCompletionRef(ctx, vn.sgClient.Client)
	if err != nil {
		return err
	}

	sg, err := sgf.Result(vn.sgClient)
	if err != nil {
		return err
	}

	rtf, err := vn.rtClient.CreateOrUpdate(ctx, vn.config.ResourceGroup, vn.names.RouteTable(), network.RouteTable{
		Location: &vn.config.Location,
	})
	if err != nil {
		return err
	}

	err = rtf.WaitForCompletionRef(ctx, vn.rtClient.Client)
	if err != nil {
		return err
	}

	rt, err := rtf.Result(vn.rtClient)
	if err != nil {
		return err
	}

	f, err := vn.vnetClient.CreateOrUpdate(ctx, vn.config.ResourceGroup, vn.names.VirtualNetwork(),
		network.VirtualNetwork{
			Location: &vn.config.Location,
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{net.CIDR},
				},
				Subnets: &[]network.Subnet{
					{
						Name: to.StringPtr(vn.names.Subnet()),
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

	err = f.WaitForCompletionRef(ctx, vn.vnetClient.Client)
	if err != nil {
		return err
	}
	return err
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

func (p *Provider) GetVirtualNetwork(ctx context.Context, net *v1alpha1.Network) (exists bool, err error) {
	_, err = p.client.vnet.Get(ctx, p.config.ResourceGroup, p.names.VirtualNetwork(), "")
	if err != nil && IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (p *Provider) EnsureVirtualNetwork(ctx context.Context, net *v1alpha1.Network) error {
	return p.create(ctx, net)
}

func (p *Provider) UpdateVirtualNetwork(ctx context.Context, net *v1alpha1.Network) error {
	vnet, err := p.client.vnet.Get(ctx, p.config.ResourceGroup, p.names.VirtualNetwork(), "")
	if err != nil {
		return err
	}

	vnet.Location = &p.config.Location
	// TODO: probably dangerous without nil checks
	(*vnet.AddressSpace.AddressPrefixes)[0] = net.CIDR
	(*vnet.AddressSpace.AddressPrefixes)[0] = net.CIDR
	(*vnet.Subnets)[0].Name = to.StringPtr(p.names.Subnet())
	(*vnet.Subnets)[0].AddressPrefix = &net.SubnetCIDR
	_, err = p.client.vnet.CreateOrUpdate(ctx, p.config.ResourceGroup, p.names.VirtualNetwork(), vnet)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureVirtualNetworkDeleted(ctx context.Context, net *v1alpha1.Network) error {
	f, err := p.client.vnet.Delete(ctx, p.config.ResourceGroup, p.names.VirtualNetwork())
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, p.client.vnet.Client)
	if err != nil {
		return err
	}

	_, err = f.Result(p.client.vnet)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) create(ctx context.Context, net *v1alpha1.Network) error {
	sgf, err := p.client.nsg.CreateOrUpdate(ctx, p.config.ResourceGroup, p.names.NetworkSecurityGroup(), network.SecurityGroup{
		Location: &p.config.Location,
	})
	if err != nil {
		return err
	}

	err = sgf.WaitForCompletionRef(ctx, p.client.nsg.Client)
	if err != nil {
		return err
	}

	sg, err := sgf.Result(p.client.nsg)
	if err != nil {
		return err
	}

	rtf, err := p.client.rt.CreateOrUpdate(ctx, p.config.ResourceGroup, p.names.RouteTable(), network.RouteTable{
		Location: &p.config.Location,
	})
	if err != nil {
		return err
	}

	err = rtf.WaitForCompletionRef(ctx, p.client.rt.Client)
	if err != nil {
		return err
	}

	rt, err := rtf.Result(p.client.rt)
	if err != nil {
		return err
	}

	f, err := p.client.vnet.CreateOrUpdate(ctx, p.config.ResourceGroup, p.names.VirtualNetwork(),
		network.VirtualNetwork{
			Location: &p.config.Location,
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{net.CIDR},
				},
				Subnets: &[]network.Subnet{
					{
						Name: to.StringPtr(p.names.Subnet()),
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

	err = f.WaitForCompletionRef(ctx, p.client.vnet.Client)
	if err != nil {
		return err
	}
	return err
}
