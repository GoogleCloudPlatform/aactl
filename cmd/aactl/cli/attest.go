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
	"github.com/GoogleCloudPlatform/aactl/pkg/attestation"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	attestCmd = &c.Command{
		Name:    "attestation",
		Aliases: []string{"att", "attest"},
		Usage:   "import attestation metadata",
		Action:  attestationCmd,
		Flags: []c.Flag{
			projectFlag,
			sourceFlag,
		},
	}
)

func attestationCmd(c *c.Context) error {
	uri, err := utils.ResolveImageURI(c.Context, c.String(sourceFlag.Name))
	if err != nil {
		return errors.Wrap(err, "error resolving source name to URI")
	}

	opt := &types.AttestationOptions{
		Project: c.String(projectFlag.Name),
		Source:  uri,
		Quiet:   isQuiet(c),
	}

	printVersion(c)

	if err := attestation.Import(c.Context, opt); err != nil {
		return errors.Wrap(err, "error executing command")
	}

	return nil
}
