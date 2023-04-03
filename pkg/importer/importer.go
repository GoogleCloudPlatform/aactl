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

package importer

import (
	"context"

	"github.com/GoogleCloudPlatform/aactl/pkg/attestation"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/vul"
	"github.com/pkg/errors"
)

type Importer func(ctx context.Context, options types.Options) error

func GetImporter(format types.SourceFormat) (Importer, error) {
	switch format {
	case types.SourceFormatSnykJSON:
		return vul.Import, nil
	case types.SourceFormatTrivyJSON:
		return attestation.Import, nil
	default:
		return nil, errors.Errorf("unimplemented importer format: %s", format)
	}
}
