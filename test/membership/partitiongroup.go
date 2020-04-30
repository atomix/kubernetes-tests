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
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

// TestPartitionGroup : integration test
func (s *TestSuite) TestPartitionGroup(t *testing.T) {
	address, err := s.getControllerAddress()
	assert.NoError(t, err)
	client, err := atomix.New(
		address,
		atomix.WithNamespace(helm.Namespace()),
		atomix.WithScope(t.Name()))
	assert.NoError(t, err)
	assert.NotNil(t, client)

	group, err := client.GetPartitionGroup(context.Background(), "test-partition-group", atomix.WithPartitions(3))
	assert.NoError(t, err)

	partition1 := group.Partition(1)
	watchCh1 := make(chan atomix.Membership)
	err = partition1.MembershipGroup().Watch(context.Background(), watchCh1)
	assert.NoError(t, err)

	partition2 := group.Partition(2)
	watchCh2 := make(chan atomix.Membership)
	err = partition2.MembershipGroup().Watch(context.Background(), watchCh2)
	assert.NoError(t, err)

	partition3 := group.Partition(3)
	watchCh3 := make(chan atomix.Membership)
	err = partition3.MembershipGroup().Watch(context.Background(), watchCh3)
	assert.NoError(t, err)

	kube, err := kubernetes.New()
	assert.NoError(t, err)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kube.Namespace(),
			Name:      "test-partition-group-member-1",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-partition-group-member",
					Image:           "atomix/test-partition-group-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-partition-group-member-1",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						"--group=test-partition-group",
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
			Name:      "test-partition-group-member-2",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-partition-group-member",
					Image:           "atomix/test-partition-group-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-partition-group-member-2",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						"--group=test-partition-group",
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
			Name:      "test-partition-group-member-3",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-partition-group-member",
					Image:           "atomix/test-partition-group-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-partition-group-member-3",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						"--group=test-partition-group",
						"--partitions=3",
						fmt.Sprintf("--test=%s", t.Name()),
					},
				},
			},
		},
	}
	_, err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Create(pod)
	assert.NoError(t, err)

	var watchJoin = func(join ...atomix.MemberID) {
		partitions := make(map[atomix.PartitionID]bool)
		for {
			select {
			case membership := <-watchCh1:
				members := make(map[atomix.MemberID]bool)
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
				members := make(map[atomix.MemberID]bool)
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
				members := make(map[atomix.MemberID]bool)
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

	members := []atomix.MemberID{"test-partition-group-member-1", "test-partition-group-member-2", "test-partition-group-member-3"}
	watchJoin(members...)

	membership := partition1.MembershipGroup().Membership()
	assert.NotNil(t, membership.Leadership)
	leader := membership.Leadership.Leader

	var gracePeriod int64
	propagation := metav1.DeletePropagationForeground
	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete(string(leader), &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod, PropagationPolicy: &propagation})
	assert.NoError(t, err)

	var watchLeave = func(leave ...atomix.MemberID) {
		partitions := make(map[atomix.PartitionID]bool)
		for {
			select {
			case membership := <-watchCh1:
				members := make(map[atomix.MemberID]bool)
				for _, member := range membership.Members {
					members[member.ID] = true
				}
				println(fmt.Sprintf("1: %v", members))
				left := true
				for _, member := range leave {
					if members[member] {
						left = false
					}
				}
				if left {
					partitions[partition1.ID] = true
					if len(partitions) == 3 {
						return
					}
				}
			case membership := <-watchCh2:
				members := make(map[atomix.MemberID]bool)
				for _, member := range membership.Members {
					members[member.ID] = true
				}
				println(fmt.Sprintf("2: %v", members))
				left := true
				for _, member := range leave {
					if members[member] {
						left = false
					}
				}
				if left {
					partitions[partition2.ID] = true
					if len(partitions) == 3 {
						return
					}
				}
			case membership := <-watchCh3:
				members := make(map[atomix.MemberID]bool)
				for _, member := range membership.Members {
					members[member.ID] = true
				}
				println(fmt.Sprintf("3: %v", members))
				left := true
				for _, member := range leave {
					if members[member] {
						left = false
					}
				}
				if left {
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

	watchLeave(leader)

	for _, member := range members {
		if member != leader {
			err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete(string(member), &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod, PropagationPolicy: &propagation})
			assert.NoError(t, err)
		}
	}

	watchLeave(members...)

	err = group.Close()
	assert.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}
