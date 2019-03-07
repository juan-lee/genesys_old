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

package controlplane

// import (
// 	"context"
// 	"time"

// 	"github.com/go-logr/logr"
// 	"github.com/juan-lee/genesys/pkg/actuator/provider"
// 	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
// 	"sigs.k8s.io/controller-runtime/pkg/reconcile"
// )

// type SingleInstance struct {
// 	log logr.Logger
// 	cpe provider.ControlPlaneEndpointFactory
// }

// func ProvideSingleInstance(log logr.Logger, cpe provider.ControlPlaneEndpointFactory) (*SingleInstance, error) {
// 	return &SingleInstance{
// 		log: log,
// 		cpe: cpe,
// 	}, nil
// }

// func (r *SingleInstance) Ensure(ctx context.Context, cp *v1alpha1.ControlPlane) (reconcile.Result, error) {
// 	r.log.Info("controlplane.EnsureDeleted enter")
// 	defer r.log.Info("controlplane.EnsureDeleted exit")

// 	if err := validate(cp); err != nil {
// 		return reconcile.Result{}, err
// 	}

// 	rec, err := r.cpe.Get(ctx, cp)
// 	if err != nil {
// 		return reconcile.Result{}, err
// 	}

// 	if err := rec.Ensure(ctx); err != nil {
// 		switch err.(type) {
// 		case *provider.ProvisioningInProgress:
// 			return reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
// 		}
// 		return reconcile.Result{}, err
// 	}
// 	return reconcile.Result{}, nil
// }

// func (r *SingleInstance) EnsureDeleted(ctx context.Context, cp *v1alpha1.ControlPlane) (reconcile.Result, error) {
// 	r.log.Info("controlplane.EnsureDeleted enter")
// 	defer r.log.Info("controlplane.EnsureDeleted exit")

// 	rec, err := r.cpe.Get(ctx, cp)
// 	if err != nil {
// 		return reconcile.Result{}, err
// 	}

// 	if err := rec.EnsureDeleted(ctx); err != nil {
// 		switch err.(type) {
// 		case *provider.ProvisioningInProgress:
// 			return reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
// 		}
// 		return reconcile.Result{}, err
// 	}
// 	return reconcile.Result{}, nil
// }

// func (r *SingleInstance) ensureAll(ctx context.Context, cp *v1alpha1.ControlPlane) error {
// 	if err := r.ensureControlPlaneEndpoint(ctx, cp); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *SingleInstance) ensureControlPlaneEndpoint(ctx context.Context, cp *v1alpha1.ControlPlane) error {
// 	return nil
// }

// func validate(cp *v1alpha1.ControlPlane) error {
// 	return nil
// }
