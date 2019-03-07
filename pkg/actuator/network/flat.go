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

package network

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

type Flat struct {
	log      logr.Logger
	provider provider.Interface
}

func ProvideFlat(log logr.Logger, cloud *v1alpha1.Cloud, provider provider.Interface) (*Flat, error) {
	return &Flat{log: log, provider: provider}, nil
}

func (b *Flat) Ensure(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.ensureVNet(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *Flat) EnsureDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.ensureVNetDeleted(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *Flat) ensureVNet(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if vnet, exists := b.provider.VirtualNetwork(); exists {
		if hasVNet, err := vnet.GetVirtualNetwork(ctx, &cluster.Spec.Network); err != nil && hasVNet {
			if err := vnet.UpdateVirtualNetwork(ctx, &cluster.Spec.Network); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		err := vnet.EnsureVirtualNetwork(ctx, &cluster.Spec.Network)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (b *Flat) ensureVNetDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if vnet, exists := b.provider.VirtualNetwork(); exists {
		if hasVNet, err := vnet.GetVirtualNetwork(ctx, &cluster.Spec.Network); err != nil && hasVNet {
			if err := vnet.EnsureVirtualNetworkDeleted(ctx, &cluster.Spec.Network); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		return nil
	}
	return nil
}
