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
	"fmt"
)

type names struct {
	prefix string
}

func (n names) VirtualNetwork() string {
	return fmt.Sprintf("%s-vnet", n.prefix)
}

func (n names) Subnet() string {
	return fmt.Sprintf("%s-subnet", n.prefix)
}

func (n names) NetworkSecurityGroup() string {
	return fmt.Sprintf("%s-nsg", n.prefix)
}

func (n names) RouteTable() string {
	return fmt.Sprintf("%s-rt", n.prefix)
}

func (n names) ControlPlaneEndpoint() string {
	return fmt.Sprintf("%s-cp-ip", n.prefix)
}
