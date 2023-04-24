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
	"context"
	"fmt"
	"regexp"
	"strings"

	ar "cloud.google.com/go/artifactregistry/apiv1"
	arpb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/pkg/errors"
)

var errStr = `Not a valid source.

A valid source can be referenced by tag or digest, has the format of
LOCATION-docker.pkg.dev/PROJECT-ID/REPOSITORY-ID/IMAGE:tag
LOCATION-docker.pkg.dev/PROJECT-ID/REPOSITORY-ID/IMAGE@sha256:digest
`

var tagRegex = regexp.MustCompile("^(?P<location>.*)-docker.pkg.dev/(?P<project>[^/]+)/(?P<repo>[^/]+)/(?P<img>.*):(?P<tag>.*)")

func ResolveImageURI(ctx context.Context, in string) (string, error) {
	// If the image already has a digest there is no need to resolve.
	if strings.Contains(in, "@sha256") {
		return in, nil
	}

	matches := tagRegex.FindAllStringSubmatch(in, -1)

	if len(matches) == 0 {
		return "", errors.Wrap(types.ErrInvalidSource, errStr)
	}

	if len(matches[0]) != 6 {
		return "", errors.Wrap(types.ErrInvalidSource, errStr)
	}

	project, location, repo, image, tag := matches[0][2], matches[0][1], matches[0][3], matches[0][4], matches[0][5]

	// Resolve to digest.
	c, err := ar.NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error creating artifact registry client")
	}
	defer c.Close()

	req := &arpb.GetTagRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/repositories/%s/packages/%s/tags/%s", project, location, repo, image, tag),
	}

	// The version returned is in the format:
	// "projects/<project>/locations/<location>/repositories/<repo-name>/packages/<imageName>/versions/sha256:<digest>"
	t, err := c.GetTag(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "failed fetching version using tag")
	}

	versionSplit := strings.Split(t.GetVersion(), "versions/")
	if len(versionSplit) != 2 {
		return "", errors.New("invalid version format returned from Artifact registry")
	}
	digest := versionSplit[1]

	return fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s@%s", location, project, repo, image, digest), nil
}
