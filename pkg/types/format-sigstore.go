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
	SigstoreFormatUnknown    SigstoreFormat = iota
	SigstoreFormatProvenance                // SLSA Provenance format

	SigstoreFormatUnknownName    = "unknown"
	SigstoreFormatProvenanceName = "provenance"
)

// SigstoreFormat represents the metatdata format.
type SigstoreFormat int64

// String returns the string representation of the format.
func (f SigstoreFormat) String() string {
	switch f {
	case SigstoreFormatProvenance:
		return SigstoreFormatProvenanceName
	default:
		return SigstoreFormatUnknownName
	}
}

// ParseSigstoreFormat parses the format.
func ParseSigstoreFormat(s string) (SigstoreFormat, error) {
	switch s {
	case SigstoreFormatProvenanceName:
		return SigstoreFormatProvenance, nil
	default:
		return SigstoreFormatUnknown, fmt.Errorf("unknown format: %s", s)
	}
}

// GetSigstoreFormats returns the supported formats.
func GetSigstoreFormats() []SigstoreFormat {
	return []SigstoreFormat{
		SigstoreFormatProvenance,
	}
}

// GetSigstoreFormatNames returns the names of the supported formats.
func GetSigstoreFormatNames() []string {
	return []string{
		SigstoreFormatProvenanceName,
	}
}
