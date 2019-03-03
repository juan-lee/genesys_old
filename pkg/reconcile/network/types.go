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

	k8sv1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler interface {
	Reconcile(k8sv1alpha1.Network) (reconcile.Result, error)
}

// Func is a function that implements the reconcile interface.
type Func func(k8sv1alpha1.Network) (reconcile.Result, error)

var _ Reconciler = Func(nil)

// Reconcile implements Reconciler.
func (r Func) Reconcile(o k8sv1alpha1.Network) (reconcile.Result, error) { return r(o) }

type VNETOptions struct {
	CIDR       string
	SubnetCIDR string
}

type VNETProvider interface {
	State(ctx context.Context, desired k8sv1alpha1.Network) error
	Update(ctx context.Context, desired k8sv1alpha1.Network) error
	Delete(ctx context.Context) error
}
