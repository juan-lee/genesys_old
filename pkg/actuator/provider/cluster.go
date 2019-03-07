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

type Status string

const (
	NeedsUpdate  Status = "NeedsUpdate"
	Succeeded    Status = "Succeeded"
	Provisioning Status = "Provisioning"
	Deleting     Status = "Deleting"
	Deleted      Status = "Deleted"
	Unknown      Status = "Unknown"
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

type ControlPlaneEndpointFactory interface {
	Get(ctx context.Context, cp *v1alpha1.ControlPlane) (Reconciler, error)
}

type VirtualNetworkFactory interface {
	Get(ctx context.Context, net *v1alpha1.Network) (Reconciler, error)
}

type NetworkSecurityGroup interface {
	Ensure(ctx context.Context, net v1alpha1.Network) error
	EnsureDeleted(ctx context.Context, net v1alpha1.Network) error
}

type VirtualNetwork interface {
	GetVirtualNetwork(ctx context.Context, net *v1alpha1.Network) (exists bool, err error)
	EnsureVirtualNetwork(ctx context.Context, net *v1alpha1.Network) error
	UpdateVirtualNetwork(ctx context.Context, net *v1alpha1.Network) error
	EnsureVirtualNetworkDeleted(ctx context.Context, net *v1alpha1.Network) error
}

type Interface interface {
	VirtualNetwork() (VirtualNetwork, bool)
}
