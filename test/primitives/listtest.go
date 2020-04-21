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
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAtomixList : integration test
func (s *TestSuite) TestAtomixList(t *testing.T) {
	client, err := s.getClient(t)
	assert.NoError(t, err)

	database, err := client.GetDatabase(context.Background(), "raft-database")
	assert.NoError(t, err)

	list, err := database.GetList(context.Background(), "TestAtomixList")
	assert.NoError(t, err)

	size, err := list.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = list.Append(context.Background(), []byte("Hello world!"))
	assert.NoError(t, err)

	size, err = list.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, size)

	value, err := list.Get(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!", string(value))

	value, err = list.Remove(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!", string(value))

	size, err = list.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = list.Append(context.Background(), []byte("Hello world!"))
	assert.NoError(t, err)

	err = list.Append(context.Background(), []byte("Hello world again!"))
	assert.NoError(t, err)

	ch := make(chan []byte)
	err = list.Items(context.Background(), ch)
	i := 0
	for value := range ch {
		if i == 0 {
			assert.Equal(t, "Hello world!", string(value))
			i++
		} else if i == 1 {
			assert.Equal(t, "Hello world again!", string(value))
			i++
		} else {
			assert.Fail(t, "Too many values")
		}
	}
	assert.NoError(t, err)
}
