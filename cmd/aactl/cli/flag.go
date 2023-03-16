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

package cli

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	c "github.com/urfave/cli/v2"
)

var (
	projectFlag = &c.StringFlag{
		Name:    "project",
		Aliases: []string{"p"},
		Usage:   "GCP project ID where the vulnerabilities will be imported",
	}

	sourceFlag = &c.StringFlag{
		Name:    "source",
		Aliases: []string{"s"},
		Usage:   "uri of the image from which the report was generated (e.g. us-docker.pkg.dev/project/repo/img@sha256:f6efe...)",
	}

	fileFlag = &c.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "path to vulnerability report to import",
	}

	formatFlag = &c.StringFlag{
		Name:    "format",
		Aliases: []string{"t"},
		Usage:   fmt.Sprintf("file type (e.g. %s, etc.)", strings.Join(types.GetSourceFormatNames(), ", ")),
	}
)
