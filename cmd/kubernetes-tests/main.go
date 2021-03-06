// Copyright 2019-present Open Networking Foundation.
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

package main

import (
	"github.com/atomix/kubernetes-tests/test/peers"
	"github.com/atomix/kubernetes-tests/test/primitives"
	"github.com/onosproject/helmit/pkg/registry"
	"github.com/onosproject/helmit/pkg/test"
)

func main() {
	registry.RegisterTestSuite("peers", &peers.TestSuite{})
	registry.RegisterTestSuite("primitives", &primitives.TestSuite{})
	test.Main()
}
