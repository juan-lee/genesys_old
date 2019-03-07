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
	"github.com/Azure/go-autorest/autorest/to"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

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

	if !reflect.DeepEqual(*net, convert(&vnet)) {
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
