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
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log  logr.Logger
	vnet provider.VirtualNetwork
	cpe  provider.ControlPlaneEndpoint
}

func ProvideReconciler(log logr.Logger, vnet provider.VirtualNetwork, cpe provider.ControlPlaneEndpoint) (*Reconciler, error) {
	return &Reconciler{
		log:  log,
		vnet: vnet,
	}, nil
}

func (r *Reconciler) Reconcile(net v1alpha1.Network) (reconcile.Result, error) {
	r.log.Info("network.Reconcile enter")
	defer r.log.Info("network.Reconcile exit")

	if result, err := r.ensureVirtualNetwork(net); err != nil {
		return result, err
	}

	if result, err := r.ensureControlPlaneEndpoint(); err != nil {
		return result, err
	}
	return reconcile.Result{}, nil
}

func (r *Reconciler) ensureVirtualNetwork(net v1alpha1.Network) (reconcile.Result, error) {
	if err := validateVirtualNetwork(net); err != nil {
		return reconcile.Result{}, err
	}

	if _, err := r.vnet.Get(context.TODO()); err != nil {
		switch err {
		case provider.ErrNotFound:
			r.log.Info("Creating VirtualNetwork")
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

func (r *Reconciler) ensureControlPlaneEndpoint() (reconcile.Result, error) {
	if err := r.cpe.Get(context.TODO(), "cpe-name"); err != nil {
		switch err {
		case provider.ErrNotFound:
			r.log.Info("Creating ControlPlaneEndpoint")
			err := r.cpe.Update(context.TODO(), "cpe-name")
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
