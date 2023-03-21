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

package container

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/pkg/errors"
)

// GetFullURL gets the fully qualified container repository url
func GetFullURL(u string) (string, error) {
	tokens := strings.Split(u, "@")
	if len(tokens) == 0 {
		return "", fmt.Errorf("GetFullURL: Invalid URL (%s)", u)
	}
	tokens = strings.Split(tokens[0], ":")
	if len(tokens) == 0 {
		return "", fmt.Errorf("GetFullURL: Invalid URL (%s)", u)
	}

	digest, err := crane.Digest(u)
	if err != nil {
		return "", errors.Wrap(err, "Getting digest from url failed")
	}

	return fmt.Sprintf("%s@%s", tokens[0], digest), nil
}
