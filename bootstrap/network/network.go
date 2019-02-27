// Copyright © 2019 The Genesys Authors
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
)

// BaseNetwork provides reconcilers for various kubernetes configurations
type BaseNetwork struct {
	vnet VNetProvider
}

// ProvideBaseNetwork provides an instance of a base networking reconciler
func ProvideBaseNetwork(vnet VNetProvider) BaseNetworkProvider {
	return &BaseNetwork{
		vnet: vnet,
	}
}

// Bootstrap provisions base networking for a kubernetes cluster
func (n BaseNetwork) Bootstrap(ctx context.Context, opt *BaseNetworkOptions) error {
	err := n.validate(opt)
	if err != nil {
		return err
	}
	return n.vnet.Bootstrap(ctx, &opt.VNet)
}

func (n BaseNetwork) validate(opt *BaseNetworkOptions) error {
	if opt == nil {
		return NewInvalidArgumentError("opt", "can't be nil")
	}
	if opt.ResourceGroup == "" {
		return NewInvalidArgumentError("ResourceGroup", "can't be empty")
	}
	if opt.Location == "" {
		return NewInvalidArgumentError("Location", "can't be empty")
	}
	return nil
}
