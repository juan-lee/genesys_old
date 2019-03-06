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

package provider

import (
	"context"
	"errors"

	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type ControlPlaneEndpoint interface {
	Ensure(ctx context.Context, ep string) error
	EnsureDeleted(ctx context.Context, ep string) error
}

type NetworkSecurityGroup interface {
	Ensure(ctx context.Context, net v1alpha1.Network) error
	EnsureDeleted(ctx context.Context, net v1alpha1.Network) error
}

type VirtualNetwork interface {
	Ensure(ctx context.Context, net *v1alpha1.Network) error
	EnsureDeleted(ctx context.Context, net *v1alpha1.Network) error
}

type Status string

const (
	NeedsUpdate  Status = "NeedsUpdate"
	Succeeded    Status = "Succeeded"
	Provisioning Status = "Provisioning"
	Deleting     Status = "Deleting"
	Deleted      Status = "Deleted"
)

type Provider interface {
	Exists() bool
	Status() Status
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}

type Reconciler interface {
	Ensure(ctx context.Context) error
	EnsureDeleted(ctx context.Context) error
}

type VirtualNetworkFactory interface {
	Get(ctx context.Context, net *v1alpha1.Network) (Reconciler, error)
}
