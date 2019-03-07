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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	"github.com/prometheus/common/log"
)

type Reconciler struct {
	log      logr.Logger
	Provider provider.Provider
}

func (r *Reconciler) Ensure(ctx context.Context) error {
	log.Info("Reconciler.Ensure", "status", r.Provider.Status())
	switch r.Provider.Status() {
	case provider.Provisioning:
		return provider.Pending(fmt.Sprintf("%T", r.Provider))
	case provider.NeedsUpdate:
		err := r.Provider.Update(ctx)
		if err != nil {
			return err
		}
		return provider.Pending(fmt.Sprintf("%T", r.Provider))
	case provider.Succeeded:
		return nil
	}
	return nil
}

func (r *Reconciler) EnsureDeleted(ctx context.Context) error {
	log.Info("Reconciler.Ensure", "status", r.Provider.Status())
	switch r.Provider.Status() {
	case provider.Deleting:
		return provider.Pending(fmt.Sprintf("%T", r.Provider))
	case provider.Succeeded:
		err := r.Provider.Delete(ctx)
		if err != nil {
			return err
		}
		return provider.Pending(fmt.Sprintf("%T", r.Provider))
	}
	return nil
}
