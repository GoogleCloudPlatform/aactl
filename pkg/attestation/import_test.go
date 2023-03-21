// Copyright 2023 Google LLC
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

package attestation

import (
	"testing"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestInvalidImport(t *testing.T) {
	err := Import(context.TODO(), nil)
	assert.Error(t, err)
	err = Import(context.TODO(), &types.AttestationOptions{})
	assert.Error(t, err)
	err = Import(context.TODO(), &types.AttestationOptions{
		Source: "us-docker.pkg.dev/project/repo/img@sha256:f6efe...",
	})
	assert.Error(t, err)
}
