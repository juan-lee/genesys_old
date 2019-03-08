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

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/juan-lee/genesys/pkg/actuator/provider"
	v1alpha1 "github.com/juan-lee/genesys/pkg/apis/kubernetes/v1alpha1"
)

type SingleInstance struct {
	log      logr.Logger
	provider provider.Interface
}

func ProvideSingleInstance(log logr.Logger, cloud *v1alpha1.Cloud, provider provider.Interface) (*SingleInstance, error) {
	return &SingleInstance{log: log, provider: provider}, nil
}

func (b *SingleInstance) Ensure(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.ensureControlPlaneEndpoint(ctx, cluster)
	if err != nil {
		return err
	}

	err := b.ensureExternalLoadBalancer(ctx, cluster)
	if err != nil {
		return err
	}

	err := b.ensureInternalLoadBalancer(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *SingleInstance) EnsureDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	err := b.ensureControlPlaneEndpointDeleted(ctx, cluster)
	if err != nil {
		return err
	}

	err := b.ensureExternalLoadBalancerDeleted(ctx, cluster)
	if err != nil {
		return err
	}

	err := b.ensureInternalLoadBalancerDeleted(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (b *SingleInstance) ensureControlPlaneEndpoint(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if cpe, supported := b.provider.ControlPlaneEndpoint(); supported {
		if exists, err := cpe.GetControlPlaneEndpoint(ctx, &cluster.Spec.ControlPlane); err != nil && exists {
			if err := cpe.UpdateControlPlaneEndpoint(ctx, &cluster.Spec.ControlPlane); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		err := cpe.EnsureControlPlaneEndpoint(ctx, &cluster.Spec.ControlPlane)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (b *SingleInstance) ensureControlPlaneEndpointDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if cpe, supported := b.provider.ControlPlaneEndpoint(); supported {
		if exists, err := cpe.GetControlPlaneEndpoint(ctx, &cluster.Spec.ControlPlane); err != nil && exists {
			if err := cpe.EnsureControlPlaneEndpointDeleted(ctx, &cluster.Spec.ControlPlane); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (b *SingleInstance) ensureExternalLoadBalancer(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if cpe, supported := b.provider.ExternalLoadBalancer(); supported {
		if exists, err := cpe.GetExternalLoadBalancer(ctx, &cluster.Spec.ControlPlane); err != nil && exists {
			if err := cpe.UpdateExternalLoadBalancer(ctx, &cluster.Spec.ControlPlane); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		err := cpe.EnsureExternalLoadBalancer(ctx, &cluster.Spec.ControlPlane)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (b *SingleInstance) ensureExternalLoadBalancerDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if cpe, supported := b.provider.ExternalLoadBalancer(); supported {
		if exists, err := cpe.GetExternalLoadBalancer(ctx, &cluster.Spec.ControlPlane); err != nil && exists {
			if err := cpe.EnsureExternalLoadBalancerDeleted(ctx, &cluster.Spec.ControlPlane); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (b *SingleInstance) ensureInternalLoadBalancer(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if cpe, supported := b.provider.InternalLoadBalancer(); supported {
		if exists, err := cpe.GetInternalLoadBalancer(ctx, &cluster.Spec.ControlPlane); err != nil && exists {
			if err := cpe.UpdateInternalLoadBalancer(ctx, &cluster.Spec.ControlPlane); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		err := cpe.EnsureInternalLoadBalancer(ctx, &cluster.Spec.ControlPlane)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (b *SingleInstance) ensureInternalLoadBalancerDeleted(ctx context.Context, cluster *v1alpha1.Cluster) error {
	if cpe, supported := b.provider.InternalLoadBalancer(); supported {
		if exists, err := cpe.GetInternalLoadBalancer(ctx, &cluster.Spec.ControlPlane); err != nil && exists {
			if err := cpe.EnsureInternalLoadBalancerDeleted(ctx, &cluster.Spec.ControlPlane); err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}

		return nil
	}
	return nil
}
