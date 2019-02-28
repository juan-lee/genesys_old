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

type vnetReconciler struct {
	vnet VNetReconciler
}

func (r vnetReconciler) Reconcile(ctx context.Context, opt *VNetOptions) error {
	err := r.validate(opt)
	if err != nil {
		return err
	}
	return r.vnet.Reconcile(ctx, opt)
}

func (r vnetReconciler) validate(opt *VNetOptions) error {
	if opt.Name == "" {
		return NewInvalidArgumentError("Name", "can't be empty")
	}
	if opt.AddressSpace == "" {
		return NewInvalidArgumentError("AddressSpace", "can't be empty")
	}
	if len(opt.Subnets) <= 0 {
		return NewInvalidArgumentError("Subnets", "can't be empty")
	}
	return nil
}
