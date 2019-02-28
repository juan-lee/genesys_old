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

	"github.com/juan-lee/genesys/reconciler/cloud"
)

type reconciler struct {
	cloud *cloud.ProviderOptions
	vnet  vnetReconciler
}

// ProvideReconciler provides an instance of a base networking reconciler
func ProvideReconciler(cloud *cloud.ProviderOptions, vnet VNetReconciler) Reconciler {
	return reconciler{
		vnet: vnetReconciler{
			vnet: vnet,
		},
	}
}

// Reconcile provisions base networking for a kubernetes cluster
func (r reconciler) Reconcile(ctx context.Context, opt *Options) error {
	err := r.validate(opt)
	if err != nil {
		return err
	}
	return r.vnet.Reconcile(ctx, &opt.VNet)
}

func (r reconciler) validate(opt *Options) error {
	if opt == nil {
		return NewInvalidArgumentError("opt", "can't be nil")
	}
	return nil
}
