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
	"github.com/GoogleCloudPlatform/aactl/pkg/container"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/vul"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

var (
	impCmd = &c.Command{
		Name:    "vulnerability",
		Aliases: []string{"vul", "vuln", "vulns"},
		Usage:   "import vulnerabilities from file",
		Action:  vulnerabilityCmd,
		Flags: []c.Flag{
			projectFlag,
			sourceFlag,
			fileFlag,
			formatFlag,
		},
	}
)

var fullResourceNameResolver = func(u string) (string, error) {
	resourceURL, err := container.GetFullURL(u)
	if err != nil {
		return "", err
	}
	log.Info().Msgf("Resource URL: %s", resourceURL)
	return resourceURL, err
}

func vulnerabilityCmd(c *c.Context) error {
	f, err := types.ParseSourceFormat(c.String(formatFlag.Name))
	if err != nil {
		return errors.Wrap(err, "error parsing source format")
	}

	opt := &types.VulnerabilityOptions{
		Project: c.String(projectFlag.Name),
		Source:  c.String(sourceFlag.Name),
		File:    c.String(fileFlag.Name),
		Format:  f,
		Quiet:   isQuiet(c),
	}

	printVersion(c)

	resourceURL, err := fullResourceNameResolver(opt.Source)
	if err != nil {
		return errors.Wrap(err, "error getting full url")
	}
	opt.Source = resourceURL

	if err := vul.Import(c.Context, opt); err != nil {
		return errors.Wrap(err, "error executing command")
	}

	return nil
}
