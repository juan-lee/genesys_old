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

// var _ provider.NetworkSecurityGroup = &NetworkSecurityGroup{}

// type NetworkSecurityGroup struct {
// 	log    logr.Logger
// 	config v1alpha1.Cloud
// 	names  *names
// 	client network.SecurityGroupsClient
// }

// func ProvideNetworkSecurityGroup(log logr.Logger, a autorest.Authorizer, c v1alpha1.Cloud, n *names) (*NetworkSecurityGroup, error) {
// 	client, err := newNSGClient(c.SubscriptionID, a)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &NetworkSecurityGroup{
// 		log:    log,
// 		config: c,
// 		names:  n,
// 		client: client,
// 	}, nil
// }

// func (r *NetworkSecurityGroup) Ensure(ctx context.Context, net v1alpha1.Network) error {
// 	// TODO
// 	return nil
// }

// func (r *NetworkSecurityGroup) EnsureDeleted(ctx context.Context, net v1alpha1.Network) error {
// 	// TODO
// 	return nil
// }
