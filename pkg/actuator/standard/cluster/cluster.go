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

package cluster

import (
	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/standard/network"
	k8sv1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log     logr.Logger
	network *network.Reconciler
}

func ProvideReconciler(log logr.Logger, net *network.Reconciler) (*Reconciler, error) {
	return &Reconciler{
		log:     log,
		network: net,
	}, nil
}

func (r *Reconciler) Reconcile(desired k8sv1alpha1.Cluster) (reconcile.Result, error) {
	r.log.Info("cluster.Reconcile enter")
	defer r.log.Info("cluster.Reconcile exit")
	result, err := r.network.Reconcile(desired.Spec.Network)
	if err != nil {
		return result, err
	}
	return reconcile.Result{}, nil
}
