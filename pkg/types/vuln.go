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

const (
	TestProjectID = "test"
)

var (
	ErrMissingProject = errors.New("missing project")
	ErrMissingFormat  = errors.New("missing format")
	ErrMissingPath    = errors.New("missing path")
	ErrMissingSource  = errors.New("missing source")
	ErrInvalidSource  = errors.New("invalid source")
)

type VulnerabilityOptions struct {
	// Project is the ID of the project to import the report into.
	Project string

	// Source is the URI of the image from which the report was generated.
	Source string

	// File path to the vulnerability report to import.
	File string

	// Format of the file to import.
	Format SourceFormat

	// Quiet suppresses output
	Quiet bool
}

func (o *VulnerabilityOptions) Validate() error {
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
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	o.Source = u.String()

	if o.File == "" {
		return ErrMissingPath
	}
	if o.Format == SourceFormatUnknown {
		return ErrMissingFormat
	}
	return nil
}