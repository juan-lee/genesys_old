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

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler interface {
	// Reconciler performs a full reconciliation for the object referred to by the Request.
	// The Controller will requeue the Request to be processed again if an error is non-nil or
	// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
	Reconcile(reconcile.Request) (reconcile.Result, error)
}

type VNETOptions struct {
	CIDR       string
	SubnetCIDR string
}

type VNETProvider interface {
	State(ctx context.Context, desired VNETOptions) error
	Update(ctx context.Context, desired VNETOptions) error
	Delete(ctx context.Context) error
}
