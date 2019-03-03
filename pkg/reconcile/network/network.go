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

	"github.com/go-logr/logr"
	k8sv1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ Reconciler = &reconciler{}

type reconciler struct {
	log     logr.Logger
	network Provider
}

func ProvideReconciler(log logr.Logger, network Provider) (Reconciler, error) {
	return &reconciler{
		log:     log,
		network: network,
	}, nil
}

func (r *reconciler) Reconcile(desired k8sv1alpha1.Network) (reconcile.Result, error) {
	r.log.Info("network.Reconcile enter")
	defer r.log.Info("network.Reconcile exit")

	err := r.network.State(context.TODO(), desired)
	if err != nil {
		r.log.Info("State", "err", err)
		err := r.network.Update(context.TODO(), desired)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
