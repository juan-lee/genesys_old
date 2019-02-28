// Copyright © 2019 The Genesys Authors
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

	"github.com/juan-lee/genesys/reconciler/network"
)

type Reconciler struct {
	net network.Reconciler
}

func ProvideSelfManaged(net network.Reconciler) *Reconciler {
	return &Reconciler{
		net: net,
	}
}

func (r Reconciler) Reconcile(ctx context.Context) error {
	return r.net.Reconcile(ctx, &network.Options{})
}
