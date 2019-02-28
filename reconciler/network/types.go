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

type NetworkOptions struct {
	ResourceGroup string
	Location      string
	VNet          VNetOptions
}

type Subnet struct {
	Name         string
	AddressSpace string
}

type VNetOptions struct {
	Name         string
	AddressSpace string
	Subnets      []Subnet
}

type VNetReconciler interface {
	Reconcile(ctx context.Context, opt *VNetOptions) error
}

type Reconciler interface {
	Reconcile(ctx context.Context, opt *NetworkOptions) error
}
