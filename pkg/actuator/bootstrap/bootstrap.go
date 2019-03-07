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

package bootstrap

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/network"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

type Bootstrap struct {
	log logr.Logger
	net *network.Flat
}

func New(cloud *v1alpha1.Cloud) (*Bootstrap, error) {
	return injectBootstrap(cloud)
}

func (b *Bootstrap) Ensure(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.ensureNetwork(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bootstrap) EnsureDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.ensureNetworkDeleted(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bootstrap) ensureNetwork(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.net.Ensure(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bootstrap) ensureNetworkDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.net.EnsureDeleted(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}
