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

package peers

import (
	"context"
	"fmt"
	atomix "github.com/atomix/go-client/pkg/client"
	"github.com/atomix/go-client/pkg/client/peer"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

// TestPeers : integration test
func (s *TestSuite) TestPeers(t *testing.T) {
	address, err := s.getControllerAddress()
	assert.NoError(t, err)
	client, err := atomix.New(
		address,
		atomix.WithNamespace(helm.Namespace()),
		atomix.WithScope(t.Name()))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	watchCh := make(chan peer.Set)
	err = client.Group().Watch(context.Background(), watchCh)
	assert.NoError(t, err)

	kube, err := kubernetes.New()
	assert.NoError(t, err)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: kube.Namespace(),
			Name:      "test-peer-1",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-peer",
					Image:           "atomix/test-peer:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-peer-1",
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
			Name:      "test-peer-2",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-peer",
					Image:           "atomix/test-peer:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-peer-2",
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
			Name:      "test-peer-3",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "test-peer",
					Image:           "atomix/test-peer:latest",
					ImagePullPolicy: corev1.PullNever,
					Args: []string{
						"test-peer-3",
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
		case peers := <-watchCh:
			println(fmt.Sprintf("%v", peers))
			if peers["test-peer-1"] != nil && peers["test-peer-2"] != nil && peers["test-peer-3"] != nil {
				break watchJoin
			}
		case <-time.After(5 * time.Minute):
			t.Fail()
			return
		}
	}

	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete("test-peer-3", &metav1.DeleteOptions{})
	assert.NoError(t, err)

watchLeave:
	for {
		select {
		case peers := <-watchCh:
			println(fmt.Sprintf("%v", peers))
			if peers["test-peer-1"] != nil && peers["test-peer-2"] != nil && peers["test-peer-3"] == nil {
				break watchLeave
			}
		case <-time.After(5 * time.Minute):
			t.Fail()
			return
		}
	}

	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete("test-peer-1", &metav1.DeleteOptions{})
	assert.NoError(t, err)

	err = kube.Clientset().CoreV1().Pods(kube.Namespace()).Delete("test-peer-2", &metav1.DeleteOptions{})
	assert.NoError(t, err)

watchAllLeave:
	for {
		select {
		case peers := <-watchCh:
			println(fmt.Sprintf("%v", peers))
			if peers["test-peer-1"] == nil && peers["test-peer-2"] == nil && peers["test-peer-3"] == nil {
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
