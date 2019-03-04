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
	"errors"

	"github.com/go-logr/logr"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"github.com/juan-lee/genesys/pkg/reconcile/provider"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log  logr.Logger
	vnet provider.VirtualNetwork
}

func ProvideReconciler(log logr.Logger, vnet provider.VirtualNetwork) (*Reconciler, error) {
	return &Reconciler{
		log:  log,
		vnet: vnet,
	}, nil
}

func (r *Reconciler) Reconcile(net v1alpha1.Network) (reconcile.Result, error) {
	r.log.Info("network.Reconcile enter")
	defer r.log.Info("network.Reconcile exit")

	result, err := r.ensureVirtualNetwork(net)
	if err != nil {
		return result, err
	}
	return reconcile.Result{}, nil
}

func (r *Reconciler) ensureVirtualNetwork(net v1alpha1.Network) (reconcile.Result, error) {
	err := validateVirtualNetwork(net)
	if err != nil {
		return reconcile.Result{}, err
	}

	_, err = r.vnet.Get(context.TODO())
	if err != nil {
		switch err {
		case provider.ErrNotFound:
			r.log.Info("Creating VirtualNetwork", "vnet", err)
			err := r.vnet.Update(context.TODO(), net)
			if err != nil {
				return reconcile.Result{}, err
			}
		default:
			r.log.Info("default", "err", err)
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

func validateVirtualNetwork(net v1alpha1.Network) error {
	if net.CIDR == "" {
		return errors.New("CIDR cannot be empty")
	}

	if net.SubnetCIDR == "" {
		return errors.New("SubnetCIDR cannot be empty")
	}

	return nil
}
