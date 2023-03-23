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
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestImport(t *testing.T) {
	// test no arg import
	set := flag.NewFlagSet("", flag.ContinueOnError)
	c := cli.NewContext(newTestApp(t), set, nil)
	err := vulnerabilityCmd(c)
	assert.Error(t, err)

	// test all formats
	for _, f := range types.GetSourceFormatNames() {
		set = flag.NewFlagSet("", flag.ContinueOnError)
		set.String(projectFlag.Name, types.TestProjectID, "")
		set.String(sourceFlag.Name, "us-docker.pkg.dev/project/repo/img@sha256:f6efe...", "")
		set.String(fileFlag.Name, fmt.Sprintf("../../../examples/data/%s.json", f), "")
		set.String(formatFlag.Name, f, "")

		c = cli.NewContext(newTestApp(t), set, nil)
		err = vulnerabilityCmd(c)
		assert.NoError(t, err)
	}
}

func newTestApp(t *testing.T) *cli.App {
	app, err := newApp("v0.0.0-test", "test", time.Now().UTC().Format(time.RFC3339))
	assert.NoError(t, err)
	return app
}
