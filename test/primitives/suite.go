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

package primitives

import (
	"github.com/onosproject/onos-test/pkg/onit/cluster"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/onosproject/onos-test/pkg/test"
)

// TestSuite is a suite of tests for Atomix primitives
type TestSuite struct {
	test.Suite
}

// SetupTestSuite sets up the Atomix test suite
func (s *TestSuite) SetupTestSuite() {
	setup.Atomix()
	protocol := cluster.GetArg("protocol").String("raft")
	switch protocol {
	case "raft":
		setup.Database("protocol").Raft()
	case "nopaxos":
		setup.Database("protocol").NOPaxos()
	default:
		setup.Database("protocol").Raft()
	}
	setup.SetupOrDie()
}
