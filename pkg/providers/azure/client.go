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
	"github.com/Azure/go-autorest/autorest"
)

const (
	userAgent string = "genesys"
)

func newVNETClient(subID string, a autorest.Authorizer) (network.VirtualNetworksClient, error) {
	client := network.NewVirtualNetworksClient(subID)
	client.Authorizer = a
	client.AddToUserAgent(userAgent)
	return client, nil
}

func newPublicIPClient(subID string, a autorest.Authorizer) (network.PublicIPAddressesClient, error) {
	client := network.NewPublicIPAddressesClient(subID)
	client.Authorizer = a
	client.AddToUserAgent(userAgent)
	return client, nil
}

func newNSGClient(subID string, a autorest.Authorizer) (network.SecurityGroupsClient, error) {
	client := network.NewSecurityGroupsClient(subID)
	client.Authorizer = a
	client.AddToUserAgent(userAgent)
	return client, nil
}

func newRouteTableClient(subID string, a autorest.Authorizer) (network.RouteTablesClient, error) {
	client := network.NewRouteTablesClient(subID)
	client.Authorizer = a
	client.AddToUserAgent(userAgent)
	return client, nil
}
