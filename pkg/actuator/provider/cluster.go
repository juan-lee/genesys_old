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

type Names interface {
	VirtualNetwork() string
}

type ControlPlaneEndpoint interface {
	Get(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
}

type SecurityGroupRule struct {
}

type NetworkSecurityGroupOptions struct {
}

type NetworkSecurityGroup interface {
	Get(ctx context.Context) error
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}

type VirtualNetwork interface {
	Get(ctx context.Context) (*v1alpha1.Network, error)
	Update(ctx context.Context, net v1alpha1.Network) error
	Delete(ctx context.Context) error
}
