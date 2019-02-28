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

package network_test

import (
	"context"
	"testing"

	"github.com/juan-lee/genesys/reconciler/network"
)

func TestValidParameters(t *testing.T) {
	var variations = []struct {
		opt      *network.Options
		expected error
	}{
		{nil, network.NewInvalidArgumentError("opt", "can't be nil")},
		{&network.Options{}, network.NewInvalidArgumentError("ResourceGroup", "can't be empty")},
		{&network.Options{ResourceGroup: "rg"}, network.NewInvalidArgumentError("Location", "can't be empty")},
		{&network.Options{Location: "westus"}, network.NewInvalidArgumentError("ResourceGroup", "can't be empty")},
		{newValidOptions(&network.VNetOptions{}), network.NewInvalidArgumentError("Name", "can't be empty")},
		{newValidOptions(&network.VNetOptions{Name: "vnet"}), network.NewInvalidArgumentError("AddressSpace", "can't be empty")},
		{newValidOptions(&network.VNetOptions{AddressSpace: "192.168.0.0/16"}), network.NewInvalidArgumentError("Name", "can't be empty")},
		{newValidOptions(&network.VNetOptions{
			Name:         "vnet",
			AddressSpace: "192.168.0.0/16",
		}), network.NewInvalidArgumentError("Subnets", "can't be empty")},
		{newValidOptions(&network.VNetOptions{
			Name:         "vnet",
			AddressSpace: "192.168.0.0/16",
			Subnets: []network.Subnet{
				{Name: "subnet", AddressSpace: "192.168.1.0/16"},
			},
		}), nil},
	}

	for _, v := range variations {
		t.Run("", func(t *testing.T) {
			net, err := network.InjectFakeReconciler(context.TODO())
			if err != nil {
				t.Error(err)
			}

			if net == nil {
				t.Errorf("expected net to be non-nil")
			}

			err = net.Reconcile(context.TODO(), v.opt)
			if err != nil && v.expected != nil && err.Error() != v.expected.Error() {
				t.Errorf("expected [%s] : actual [%s]", v.expected, err.Error())
			}
			if (err == nil && v.expected != nil) || (err != nil && v.expected == nil) {
				t.Errorf("expected [%v] : actual [%v]", v.expected, err)
			}
		})
	}
}

func newValidOptions(vnet *network.VNetOptions) *network.Options {
	return &network.Options{
		ResourceGroup: "valid_resource_group",
		Location:      "valid_location",
		VNet:          *vnet,
	}
}
