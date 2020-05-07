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
	"github.com/atomix/go-client/pkg/client/cluster"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

// TestMembership : integration test
func (s *TestSuite) TestMembership(t *testing.T) {
	address, err := s.getControllerAddress()
	assert.NoError(t, err)
	client, err := atomix.New(
		address,
		atomix.WithNamespace(helm.Namespace()),
		atomix.WithScope(t.Name()))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	watchCh := make(chan cluster.Membership)
	err = client.Cluster().Watch(context.Background(), watchCh)
	assert.NoError(t, err)

	kube, err := kubernetes.New()
	assert.NoError(t, err)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kube.Namespace(),
			Name:      "test-member-1",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-member",
					Image:           "atomix/test-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-member-1",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
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
			Name:      "test-member-2",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-member",
					Image:           "atomix/test-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-member-2",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
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
			Name:      "test-member-3",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-member",
					Image:           "atomix/test-member:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-member-3",
						fmt.Sprintf("--controller=%s", address),
						fmt.Sprintf("--namespace=%s", kube.Namespace()),
						fmt.Sprintf("--test=%s", t.Name()),
					},
				},
			},
		},
	}
	_, err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Create(pod)
	assert.NoError(t, err)

watchJoin:
	for {
		select {
		case membership := <-watchCh:
			members := make(map[cluster.MemberID]bool)
			for _, member := range membership.Members {
				members[member.ID] = true
			}
			println(fmt.Sprintf("%v", members))
			if members["test-member-1"] && members["test-member-2"] && members["test-member-3"] {
				break watchJoin
			}
		case <-time.After(5 * time.Minute):
			t.Fail()
			return
		}
	}

	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete("test-member-3", &metav1.DeleteOptions{})
	assert.NoError(t, err)

watchLeave:
	for {
		select {
		case membership := <-watchCh:
			members := make(map[cluster.MemberID]bool)
			for _, member := range membership.Members {
				members[member.ID] = true
			}
			println(fmt.Sprintf("%v", members))
			if members["test-member-1"] && members["test-member-2"] && !members["test-member-3"] {
				break watchLeave
			}
		case <-time.After(5 * time.Minute):
			t.Fail()
			return
		}
	}

	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete("test-member-1", &metav1.DeleteOptions{})
	assert.NoError(t, err)

	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete("test-member-2", &metav1.DeleteOptions{})
	assert.NoError(t, err)

watchAllLeave:
	for {
		select {
		case membership := <-watchCh:
			members := make(map[cluster.MemberID]bool)
			for _, member := range membership.Members {
				members[member.ID] = true
			}
			println(fmt.Sprintf("%v", members))
			if !members["test-member-1"] && !members["test-member-2"] && !members["test-member-3"] {
				break watchAllLeave
			}
		case <-time.After(5 * time.Minute):
			t.Fail()
			return
		}
	}

	err = client.Close()
	assert.NoError(t, err)
}
