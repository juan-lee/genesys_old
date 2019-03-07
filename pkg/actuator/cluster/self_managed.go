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
	"context"

	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/controlplane"
	"github.com/juan-lee/genesys/pkg/actuator/network"
	k8sv1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SelfManaged struct {
	log     logr.Logger
	network *network.Flat
	cp      *controlplane.SingleInstance
}

func ProvideSelfManaged(log logr.Logger, net *network.Flat, cp *controlplane.SingleInstance) (*SelfManaged, error) {
	return &SelfManaged{
		log:     log,
		network: net,
	}, nil
}

func (r *SelfManaged) Ensure(ctx context.Context, desired k8sv1alpha1.Cluster) (reconcile.Result, error) {
	r.log.Info("cluster.Update enter")
	defer r.log.Info("cluster.Update exit")
	result, err := r.network.Ensure(ctx, &desired.Spec.Network)
	if err != nil {
		return result, err
	}

	result, err = r.cp.Ensure(ctx, &desired.Spec.ControlPlane)
	if err != nil {
		return result, err
	}
	return reconcile.Result{}, nil
}

func (r *SelfManaged) EnsureDeleted(ctx context.Context, desired k8sv1alpha1.Cluster) (reconcile.Result, error) {
	// TODO
	return reconcile.Result{}, nil
}
