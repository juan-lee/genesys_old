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
	"time"

	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Flat struct {
	log  logr.Logger
	vnet provider.VirtualNetworkFactory
}

func ProvideFlatNetwork(log logr.Logger, vnet provider.VirtualNetworkFactory) (*Flat, error) {
	return &Flat{
		log:  log,
		vnet: vnet,
	}, nil
}

func (r *Flat) Ensure(ctx context.Context, net *v1alpha1.Network) (reconcile.Result, error) {
	r.log.Info("network.Update enter")
	defer r.log.Info("network.Update exit")

	if err := validateVirtualNetwork(net); err != nil {
		return reconcile.Result{}, err
	}
	rec, err := r.vnet.Get(ctx, net)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := rec.Ensure(ctx); err != nil {
		switch err.(type) {
		case *provider.ProvisioningInProgress:
			return reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
		}
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *Flat) EnsureDeleted(ctx context.Context, net *v1alpha1.Network) (reconcile.Result, error) {
	// TODO
	return reconcile.Result{}, nil
}

func validateVirtualNetwork(net *v1alpha1.Network) error {
	if net.CIDR == "" {
		return errors.New("CIDR cannot be empty")
	}

	if net.SubnetCIDR == "" {
		return errors.New("SubnetCIDR cannot be empty")
	}

	return nil
}
