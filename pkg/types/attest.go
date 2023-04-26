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

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// AttestationOptions are the options for importing an attestation.
type AttestationOptions struct {
	// Project is the ID of the project to import the report into.
	Project string

	// Source is the URI of the image from which the report was generated.
	Source string

	// Quiet suppresses output
	Quiet bool
}

// Validate validates the options.
func (o *AttestationOptions) Validate() error {
	if o.Project == "" {
		return ErrMissingProject
	}

	// Validate URL and ensure that scheme is specified
	if o.Source == "" {
		return ErrMissingSource
	}

	u, err := url.Parse(o.Source)
	if err != nil {
		return errors.Wrap(ErrInvalidSource, err.Error())
	}
	u.Scheme = ""
	uri, _ := strings.CutPrefix(u.String(), "//")
	o.Source = uri

	return nil
}
