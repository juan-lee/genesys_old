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

	"github.com/juan-lee/genesys/bootstrap/network"
)

func TestValidParameters(t *testing.T) {
	var variations = []struct {
		opt      *network.BaseNetworkOptions
		expected string
	}{
		{nil, "opt is invalid : can't be nil"},
		{&network.BaseNetworkOptions{}, "ResourceGroup is invalid : can't be empty"},
		{&network.BaseNetworkOptions{ResourceGroup: "rg"}, "Location is invalid : can't be empty"},
		{&network.BaseNetworkOptions{Location: "westus"}, "ResourceGroup is invalid : can't be empty"},
	}

	for _, v := range variations {
		t.Run(v.expected, func(t *testing.T) {
			net, err := network.InjectFakeBaseNetwork(context.TODO())
			if err != nil {
				t.Error(err)
			}

			if net == nil {
				t.Errorf("expected net to be non-nil")
			}

			err = net.Bootstrap(context.TODO(), v.opt)
			if err.Error() != v.expected {
				t.Errorf("expected [%s] : actual [%s]", v.expected, err.Error())
			}
		})
	}
}
