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
	"github.com/GoogleCloudPlatform/aactl/pkg/sigstore"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	impSigstoreCmd = &c.Command{
		Name:    "import-sigstore",
		Aliases: []string{"impst"},
		Usage:   "import Sigstore-format metadata",
		Action:  importSigstoreCmd,
		Flags: []c.Flag{
			projectFlag,
			sourceFlag,
			formatFlag,
		},
	}
)

func importSigstoreCmd(c *c.Context) error {
	f, err := types.ParseSigstoreFormat(c.String(formatFlag.Name))
	if err != nil {
		return errors.Wrap(err, "error parsing source format")
	}

	opt := &types.ImportSigstoreOptions{
		Project: c.String(projectFlag.Name),
		Source:  c.String(sourceFlag.Name),
		Format:  f,
		Quiet:   isQuiet(c),
	}

	printVersion(c)

	if err := sigstore.Import(c.Context, opt); err != nil {
		return errors.Wrap(err, "error executing command")
	}

	return nil
}
