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

	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

type ControlPlaneEndpoint interface {
	GetControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) (exists bool, err error)
	EnsureControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) error
	UpdateControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) error
	EnsureControlPlaneEndpointDeleted(ctx context.Context, cp *v1alpha1.ControlPlane) error
}

type ExternalLoadBalancer interface {
	GetExternalLoadBalancer(ctx context.Context, cp *v1alpha1.ControlPlane) (exists bool, err error)
	EnsureExternalLoadBalancer(ctx context.Context, cp *v1alpha1.ControlPlane) error
	UpdateExternalLoadBalancer(ctx context.Context, cp *v1alpha1.ControlPlane) error
	EnsureExternalLoadBalancerDeleted(ctx context.Context, cp *v1alpha1.ControlPlane) error
}

type InternalLoadBalancer interface {
	GetInternalLoadBalancer(ctx context.Context, cp *v1alpha1.ControlPlane) (exists bool, err error)
	EnsureInternalLoadBalancer(ctx context.Context, cp *v1alpha1.ControlPlane) error
	UpdateInternalLoadBalancer(ctx context.Context, cp *v1alpha1.ControlPlane) error
	EnsureInternalLoadBalancerDeleted(ctx context.Context, cp *v1alpha1.ControlPlane) error
}

type VirtualNetwork interface {
	GetVirtualNetwork(ctx context.Context, net *v1alpha1.Network) (exists bool, err error)
	EnsureVirtualNetwork(ctx context.Context, net *v1alpha1.Network) error
	UpdateVirtualNetwork(ctx context.Context, net *v1alpha1.Network) error
	EnsureVirtualNetworkDeleted(ctx context.Context, net *v1alpha1.Network) error
}

type Interface interface {
	VirtualNetwork() (VirtualNetwork, bool)
	ControlPlaneEndpoint() (ControlPlaneEndpoint, bool)
	ExternalLoadBalancer() (ExternalLoadBalancer, bool)
	InternalLoadBalancer() (InternalLoadBalancer, bool)
}
