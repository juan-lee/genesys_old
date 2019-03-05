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

// import (
// 	"context"

// 	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
// 	"github.com/Azure/go-autorest/autorest"
// 	"github.com/go-logr/logr"
// 	"github.com/juan-lee/genesys/pkg/actuator/provider"
// 	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
// )

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

// func (r *NetworkSecurityGroup) Get(ctx context.Context) error {
// 	_, err := r.get(ctx)
// 	return err
// }

// func (r *NetworkSecurityGroup) Update(ctx context.Context) error {
// 	vnet, err := r.get(ctx)
// 	if err != nil && err == provider.ErrNotFound {
// 		f, err := r.client.CreateOrUpdate(ctx, r.config.ResourceGroup, r.names.NetworkSecurityGroup(), network.SecurityGroup{
// 			Location:                      &r.config.Location,
// 			SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{},
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		err = f.WaitForCompletionRef(ctx, r.client.Client)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	} else if err != nil {
// 		return err
// 	}

// 	f, err := r.client.CreateOrUpdate(ctx, r.config.ResourceGroup, r.names.NetworkSecurityGroup(), network.SecurityGroup{
// 		Location:                      &r.config.Location,
// 		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	err = f.WaitForCompletionRef(ctx, r.client.Client)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *NetworkSecurityGroup) Delete(ctx context.Context) error {
// 	return nil
// }

// func (r *NetworkSecurityGroup) get(ctx context.Context) (*network.SecurityGroup, error) {
// 	vnet, err := r.client.Get(ctx, r.config.ResourceGroup, r.names.NetworkSecurityGroup(), "")
// 	if err != nil {
// 		if derr, ok := err.(autorest.DetailedError); ok && derr.StatusCode == 404 {
// 			return nil, provider.ErrNotFound
// 		}
// 		return nil, err
// 	}

// 	return &vnet, nil
// }
