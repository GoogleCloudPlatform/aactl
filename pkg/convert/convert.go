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

package convert

import (
	"context"

	"github.com/GoogleCloudPlatform/aactl/pkg/convert/grype"
	"github.com/GoogleCloudPlatform/aactl/pkg/convert/snyk"
	"github.com/GoogleCloudPlatform/aactl/pkg/convert/trivy"
	"github.com/GoogleCloudPlatform/aactl/pkg/src"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/pkg/errors"
)

// VulnerabilityConverter is a function that converts a source to a list of AA notes.
type VulnerabilityConverter func(ctx context.Context, s *src.Source) (map[string]types.NoteOccurrences, error)

// GetConverter returns a vulnerability converter for the given format.
func GetConverter(format types.SourceFormat) (VulnerabilityConverter, error) {
	switch format {
	case types.SourceFormatSnykJSON:
		return snyk.Convert, nil
	case types.SourceFormatTrivyJSON:
		return trivy.Convert, nil
	case types.SourceFormatGrypeJSON:
		return grype.Convert, nil
	default:
		return nil, errors.Errorf("unimplemented conversion format: %s", format)
	}
}
