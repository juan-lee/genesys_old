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
	"github.com/juan-lee/genesys/pkg/actuator/network"
	k8sv1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Actuator struct {
	log     logr.Logger
	network *network.Actuator
}

func ProvideActuator(log logr.Logger, net *network.Actuator) (*Actuator, error) {
	return &Actuator{
		log:     log,
		network: net,
	}, nil
}

func (r *Actuator) Update(desired k8sv1alpha1.Cluster) (reconcile.Result, error) {
	r.log.Info("cluster.Update enter")
	defer r.log.Info("cluster.Update exit")
	result, err := r.network.Update(desired.Spec.Network)
	if err != nil {
		return result, err
	}
	return reconcile.Result{}, nil
}
