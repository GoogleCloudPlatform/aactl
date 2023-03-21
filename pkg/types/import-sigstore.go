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

	"github.com/pkg/errors"
)

type ImportSigstoreOptions struct {
	// Project is the ID of the project to import the report into.
	Project string

	// Source is the URI of the image from which the report was generated.
	Source string

	// Format of the metadata type to import.
	Format SigstoreFormat

	// Quiet suppresses output
	Quiet bool
}

func (i *ImportSigstoreOptions) Validate() error {
	if i.Project == "" {
		return ErrMissingProject
	}

	// Validate URL and ensure that scheme is specified
	if i.Source == "" {
		return ErrMissingSource
	}

	u, err := url.Parse(i.Source)
	if err != nil {
		return errors.Wrap(ErrInvalidSource, err.Error())
	}
	if u.Scheme == "" {
		u.Scheme = ""
	}
	i.Source = u.String()

	if i.Format == SigstoreFormatUnknown {
		return ErrMissingFormat
	}
	return nil
}
