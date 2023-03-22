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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestAttestationImport(t *testing.T) {
	// test no arg import
	set := flag.NewFlagSet("", flag.ContinueOnError)
	c := cli.NewContext(newTestApp(t), set, nil)
	err := attestationCmd(c)
	assert.Error(t, err)

	set = flag.NewFlagSet("", flag.ContinueOnError)
	set.String(projectFlag.Name, "cloudy-build", "")
	set.String(sourceFlag.Name, "us-west1-docker.pkg.dev/cloudy-build/demo/django@sha256:86a8fb755258259703f3a780c6502df24c98953293fe441c0035113a15730710", "")

	c = cli.NewContext(newTestApp(t), set, nil)
	err = attestationCmd(c)
	assert.NoError(t, err)
}
