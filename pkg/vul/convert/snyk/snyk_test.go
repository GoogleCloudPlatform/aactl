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

package snyk

import (
	"testing"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSnykConverter(t *testing.T) {
	opt := &types.VulnerabilityOptions{
		Project: types.TestProjectID,
		Source:  "us-docker.pkg.dev/project/repo/img@sha256:f6efe...",
		File:    "../../../../examples/data/snyk.json",
		Format:  types.SourceFormatSnykJSON,
	}
	s, err := utils.NewFileSource(opt.File, opt.Source)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	list, err := Convert(s)
	assert.NoErrorf(t, err, "failed to convert: %v", err)
	assert.NotNil(t, list)

	for id, nocc := range list {
		n := nocc.Note
		assert.NotEmpty(t, id)
		assert.NotEmpty(t, n.Name)
		assert.NotEmpty(t, n.ShortDescription)
		assert.NotEmpty(t, n.LongDescription)
		assert.NotEmpty(t, n.RelatedUrl)
		for _, u := range n.RelatedUrl {
			assert.NotEmpty(t, u.Label)
			assert.NotEmpty(t, u.Url)
		}
		assert.NotNil(t, n.GetVulnerability())
		assert.NotEmpty(t, n.GetVulnerability().CvssScore)
		assert.NotNil(t, n.GetVulnerability().CvssVersion)
		assert.NotNil(t, n.GetVulnerability().CvssV3)
		assert.NotEmpty(t, n.GetVulnerability().CvssV3.BaseScore)
		assert.NotEmpty(t, n.GetVulnerability().Severity)
		assert.NotEmpty(t, n.GetVulnerability().Details)
		for _, d := range n.GetVulnerability().Details {
			assert.NotEmpty(t, d.AffectedPackage)
			assert.NotEmpty(t, d.AffectedCpeUri)
		}
	}
}
