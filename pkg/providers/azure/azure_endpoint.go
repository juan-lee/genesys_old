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

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

var _ provider.ControlPlaneEndpoint = &ControlPlaneEndpoint{}

type ControlPlaneEndpoint struct {
	log    logr.Logger
	config v1alpha1.Cloud
	client network.PublicIPAddressesClient
}

func ProvideControlPlaneEndpoint(log logr.Logger, a autorest.Authorizer, c v1alpha1.Cloud) (*ControlPlaneEndpoint, error) {
	client, err := newPublicIPClient(c.SubscriptionID, a)
	if err != nil {
		return nil, err
	}
	return &ControlPlaneEndpoint{
		log:    log,
		config: c,
		client: client,
	}, nil
}

func (cpe *ControlPlaneEndpoint) Get(ctx context.Context, name string) error {
	_, err := cpe.client.Get(ctx, cpe.config.ResourceGroup, name, "")
	if err != nil {
		if derr, ok := err.(autorest.DetailedError); ok && derr.StatusCode == 404 {
			return provider.ErrNotFound
		}
		return err
	}
	return nil
}

func (cpe *ControlPlaneEndpoint) Update(ctx context.Context, name string) error {
	f, err := cpe.client.CreateOrUpdate(ctx, cpe.config.ResourceGroup, name, network.PublicIPAddress{
		Name:     &name,
		Location: &cpe.config.Location,
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			PublicIPAddressVersion:   network.IPv4,
			PublicIPAllocationMethod: network.Static,
		},
	})
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, cpe.client.Client)
	if err != nil {
		return err
	}
	return nil
}

func (cpe *ControlPlaneEndpoint) Delete(ctx context.Context, name string) error {
	f, err := cpe.client.Delete(ctx, cpe.config.ResourceGroup, name)
	if err != nil {
		return err
	}

	err = f.WaitForCompletionRef(ctx, cpe.client.Client)
	if err != nil {
		return err
	}
	return nil
}
