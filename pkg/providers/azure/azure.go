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
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

var _ provider.Interface = (*Provider)(nil)
var _ provider.ControlPlaneEndpoint = (*Provider)(nil)
var _ provider.VirtualNetwork = (*Provider)(nil)

type client struct {
	vnet network.VirtualNetworksClient
	nsg  network.SecurityGroupsClient
	rt   network.RouteTablesClient
	pip  network.PublicIPAddressesClient
}

type Provider struct {
	log    logr.Logger
	config *v1alpha1.Cloud
	names  *names
	client *client
}

func NewProvider(cloud *v1alpha1.Cloud) (*Provider, error) {
	return injectProvider(cloud)
}

func (p *Provider) ControlPlaneEndpoint() (provider.ControlPlaneEndpoint, bool) {
	return p, true
}

func (p *Provider) VirtualNetwork() (provider.VirtualNetwork, bool) {
	return p, true
}
