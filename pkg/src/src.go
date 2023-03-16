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

package src

import (
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
)

// NewSource returns a new Source from the given path.
func NewSource(opt *types.ImportOptions) (*Source, error) {
	if opt == nil {
		return nil, errors.New("options required")
	}

	if opt.File == "" {
		return nil, errors.New("file is required")
	}

	c, err := gabs.ParseJSONFile(opt.File)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse file: %s", opt.File)
	}

	s := &Source{
		URI:  opt.Source,
		Data: c,
	}

	return s, nil
}

type Source struct {
	// URI is the image URI.
	URI string

	// Data is the source data.
	Data *gabs.Container
}
