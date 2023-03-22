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
	"github.com/GoogleCloudPlatform/aactl/pkg/attestation/convert/provenance02"
	"github.com/GoogleCloudPlatform/aactl/pkg/provenance"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	"github.com/pkg/errors"
	g "google.golang.org/genproto/googleapis/grafeas/v1"
)

type Converter func(nr utils.NoteResource, resourceURL string, env *provenance.Envelope) (*g.Note, *g.Occurrence, error)

// GetConverter returns an attestation converter for the given format.
func GetConverter(intotoType string, intotoPredicateType string) (Converter, error) {
	if intotoType == "https://in-toto.io/Statement/v0.1" {
		switch intotoPredicateType {
		case "https://slsa.dev/provenance/v0.2":
			return provenance02.Convert, nil
		}
	}
	return nil, errors.Wrapf(types.ErrorNotSupported,
		"unimplemented env format: %s, %s", intotoType, intotoPredicateType)
}
