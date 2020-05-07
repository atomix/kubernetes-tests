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

package membership

import (
	"context"
	"fmt"
	atomix "github.com/atomix/go-client/pkg/client"
	"github.com/atomix/go-client/pkg/client/partition"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

// TestMap : integration test
func (s *TestSuite) TestMap(t *testing.T) {
	address, err := s.getControllerAddress()
	assert.NoError(t, err)
	client, err := atomix.New(
		address,
		atomix.WithNamespace(helm.Namespace()),
		atomix.WithScope(t.Name()))
	assert.NoError(t, err)
	assert.NotNil(t, client)

	group, err := client.GetPartitionGroup(context.Background(), "test-partition-group-map", partition.WithPartitions(3))
	assert.NoError(t, err)

	partition1 := group.Partition(1)
	watchCh1 := make(chan partition.Membership)
	err = partition1.Watch(context.Background(), watchCh1)
	assert.NoError(t, err)

	partition2 := group.Partition(2)
	watchCh2 := make(chan partition.Membership)
	err = partition2.Watch(context.Background(), watchCh2)
	assert.NoError(t, err)

	partition3 := group.Partition(3)
	watchCh3 := make(chan partition.Membership)
	err = partition3.Watch(context.Background(), watchCh3)
	assert.NoError(t, err)

	kube, err := kubernetes.New()
	assert.NoError(t, err)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kube.Namespace(),
			Name:      "test-partition-group-map-member-1",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-partition-group-map-member",
					Image:           "atomix/test-partition-group-map-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-partition-group-map-member-1",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						"--group=test-partition-group-map",
						"--partitions=3",
						fmt.Sprintf("--test=%s", t.Name()),
					},
				},
			},
		},
	}
	_, err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Create(pod)
	assert.NoError(t, err)

	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kube.Namespace(),
			Name:      "test-partition-group-map-member-2",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-partition-group-map-member",
					Image:           "atomix/test-partition-group-map-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-partition-group-map-member-2",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						"--group=test-partition-group-map",
						"--partitions=3",
						fmt.Sprintf("--test=%s", t.Name()),
					},
				},
			},
		},
	}
	_, err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Create(pod)
	assert.NoError(t, err)

	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kube.Namespace(),
			Name:      "test-partition-group-map-member-3",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-partition-group-map-member",
					Image:           "atomix/test-partition-group-map-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-partition-group-map-member-3",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						"--group=test-partition-group-map",
						"--partitions=3",
						fmt.Sprintf("--test=%s", t.Name()),
					},
				},
			},
		},
	}
	_, err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Create(pod)
	assert.NoError(t, err)

	var watchJoin = func(join ...partition.MemberID) {
		partitions := make(map[partition.ID]bool)
		for {
			select {
			case membership := <-watchCh1:
				members := make(map[partition.MemberID]bool)
				for _, member := range membership.Members {
					members[member.ID] = true
				}
				println(fmt.Sprintf("1: %v", members))
				if len(membership.Members) > 0 {
					assert.NotNil(t, membership.Leadership)
				}
				joined := true
				for _, member := range join {
					if !members[member] {
						joined = false
					}
				}
				if joined {
					partitions[partition1.ID] = true
					if len(partitions) == 3 {
						return
					}
				}
			case membership := <-watchCh2:
				members := make(map[partition.MemberID]bool)
				for _, member := range membership.Members {
					members[member.ID] = true
				}
				println(fmt.Sprintf("2: %v", members))
				if len(membership.Members) > 0 {
					assert.NotNil(t, membership.Leadership)
				}
				joined := true
				for _, member := range join {
					if !members[member] {
						joined = false
					}
				}
				if joined {
					partitions[partition2.ID] = true
					if len(partitions) == 3 {
						return
					}
				}
			case membership := <-watchCh3:
				members := make(map[partition.MemberID]bool)
				for _, member := range membership.Members {
					members[member.ID] = true
				}
				println(fmt.Sprintf("3: %v", members))
				if len(membership.Members) > 0 {
					assert.NotNil(t, membership.Leadership)
				}
				joined := true
				for _, member := range join {
					if !members[member] {
						joined = false
					}
				}
				if joined {
					partitions[partition3.ID] = true
					if len(partitions) == 3 {
						return
					}
				}
			case <-time.After(1 * time.Minute):
				t.Fail()
				return
			}
		}
	}

	members := []partition.MemberID{"test-partition-group-map-member-1", "test-partition-group-map-member-2", "test-partition-group-map-member-3"}
	watchJoin(members...)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	m, err := group.GetMap(ctx, t.Name())
	cancel()
	assert.NoError(t, err)
	assert.NotNil(t, m)

	entry, err := m.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.Nil(t, entry)

	entry, err = m.Put(context.Background(), "foo", []byte("bar"))
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "foo", entry.Key)
	assert.Equal(t, "bar", string(entry.Value))
	assert.NotEqual(t, 0, entry.Version)
	version := entry.Version

	entry, err = m.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "foo", entry.Key)
	assert.Equal(t, "bar", string(entry.Value))
	assert.Equal(t, version, entry.Version)

	size, err := m.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, size)

	entry, err = m.Remove(context.Background(), "foo")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "foo", entry.Key)
	assert.Equal(t, "bar", string(entry.Value))
	assert.Equal(t, version, entry.Version)

	entry, err = m.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.Nil(t, entry)

	size, err = m.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = group.Close(context.Background())
	assert.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}
