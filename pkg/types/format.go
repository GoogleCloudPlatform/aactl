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

package types

import "fmt"

const (
	SourceFormatUnknown   SourceFormat = iota
	SourceFormatGrypeJSON              // grype JSON format
	SourceFormatTrivyJSON              // trivy JSON format
	SourceFormatSnykJSON               // snyk JSON format

	SourceFormatUnknownName   = "unknown"
	SourceFormatGrypeJSONName = "grype"
	SourceFormatTrivyJSONName = "trivy"
	SourceFormatSnykJSONName  = "snyk"
)

// SourceFormat represents the source format.
type SourceFormat int64

// String returns the string representation of the source format.
func (f SourceFormat) String() string {
	switch f {
	case SourceFormatGrypeJSON:
		return SourceFormatGrypeJSONName
	case SourceFormatTrivyJSON:
		return SourceFormatTrivyJSONName
	case SourceFormatSnykJSON:
		return SourceFormatSnykJSONName
	default:
		return SourceFormatUnknownName
	}
}

// ParseSourceFormat parses the source format.
func ParseSourceFormat(s string) (SourceFormat, error) {
	switch s {
	case SourceFormatGrypeJSONName:
		return SourceFormatGrypeJSON, nil
	case SourceFormatTrivyJSONName:
		return SourceFormatTrivyJSON, nil
	case SourceFormatSnykJSONName:
		return SourceFormatSnykJSON, nil
	default:
		return SourceFormatUnknown, fmt.Errorf("unknown format: %s", s)
	}
}

// GetSourceFormats returns the supported source formats.
func GetSourceFormats() []SourceFormat {
	return []SourceFormat{
		SourceFormatGrypeJSON,
		SourceFormatTrivyJSON,
		SourceFormatSnykJSON,
	}
}

// GetSourceFormatNames returns the names of the supported source formats.
func GetSourceFormatNames() []string {
	return []string{
		SourceFormatGrypeJSONName,
		SourceFormatTrivyJSONName,
		SourceFormatSnykJSONName,
	}
}
