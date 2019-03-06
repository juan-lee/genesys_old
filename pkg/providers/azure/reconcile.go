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

package azure

import (
	"context"

	"github.com/juan-lee/genesys/pkg/actuator/provider"
)

type Reconciler struct {
	Provider provider.Provider
}

func (r *Reconciler) Ensure(ctx context.Context) error {
	switch r.Provider.Status() {
	case provider.NeedsUpdate:
		return r.Provider.Update(ctx)
	case provider.Provisioning:
		return provider.Pending("Concrete")
	case provider.Succeeded:
		return nil
	}
	return nil
}

func (r *Reconciler) EnsureDeleted(ctx context.Context) error {
	switch r.Provider.Status() {
	case provider.Deleting:
		return provider.Pending("Concrete")
	case provider.Succeeded:
		err := r.Provider.Delete(ctx)
		if err != nil {
			return err
		}
		return provider.Pending("Concrete")
	}
	return nil
}
