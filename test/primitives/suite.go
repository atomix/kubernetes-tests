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
	"errors"
	atomix "github.com/atomix/go-client/pkg/client"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/onosproject/helmit/pkg/test"
	"testing"
)

// TestSuite is a suite of tests for Atomix primitives
type TestSuite struct {
	test.Suite
}

// getClient returns the client for the test cluster
func (s *TestSuite) getClient(t *testing.T) (*atomix.Client, error) {
	kube, err := kubernetes.NewForRelease(helm.Release("atomix-controller"))
	if err != nil {
		return nil, err
	}
	services, err := kube.CoreV1().Services().List()
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, errors.New("no Atomix service found")
	}
	service := services[0]
	address := service.Ports()[0].Address(true)
	return atomix.New(
		address,
		atomix.WithNamespace(helm.Namespace()),
		atomix.WithScope(t.Name()))
}

// SetupTestSuite sets up the Atomix test suite
func (s *TestSuite) SetupTestSuite() error {
	err := helm.Chart("kubernetes-controller", "https://charts.atomix.io").
		Release("atomix-controller").
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	err = helm.Chart("raft-storage-controller", "https://charts.atomix.io").
		Release("raft-storage-controller").
		Set("scope", "Namespace").
		Install(true)
	if err != nil {
		return err
	}

	err = helm.Chart("raft-database", "https://charts.atomix.io").
		Release("raft-database").
		Set("clusters", 1).
		Set("partitions", 1).
		Set("replicas", 1).
		Install(true)
	if err != nil {
		return err
	}
	return nil
}
