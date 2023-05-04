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

package utils

import (
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
)

// NewFileSource returns a new Source from the given path.
func NewFileSource(project, path, uri string) (*Source, error) {
	if path == "" {
		return nil, types.ErrMissingPath
	}

	c, err := gabs.ParseJSONFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse file: %s", path)
	}

	f := discoverFormat(c)
	if f == types.SourceFormatUnknown {
		return nil, types.ErrInvalidFormat
	}

	s := &Source{
		Project: project,
		URI:     uri,
		Data:    c,
		Format:  f,
	}

	return s, nil
}

type Source struct {
	// Project is the associated project.
	Project string

	// URI is the image URI.
	URI string

	// Format of the file to import.
	Format types.SourceFormat

	// Data is the source data.
	Data *gabs.Container
}

func discoverFormat(c *gabs.Container) types.SourceFormat {
	if c == nil {
		return types.SourceFormatUnknown
	}

	// grype
	d := c.Search("descriptor", "name")
	if d.Exists() && d.Data() != nil && d.Data().(string) == "grype" {
		return types.SourceFormatGrypeJSON
	}

	// trivy
	if c.ExistsP("SchemaVersion") && c.ExistsP("Results") {
		return types.SourceFormatTrivyJSON
	}

	// snyk
	if c.Search("vulnerabilities").Exists() && c.Search("applications").Exists() {
		return types.SourceFormatSnykJSON
	}

	return types.SourceFormatUnknown
}
