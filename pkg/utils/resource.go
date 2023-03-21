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
	"fmt"
	"strings"
)

type NoteResource struct {
	Project string
	NoteID  string
}

func GetNoteResource(p string) (*NoteResource, error) {
	tokens := strings.Split(p, "/")
	if len(tokens) < 4 {
		return nil, fmt.Errorf("unable to parse resource (%s)", p)
	}
	project := tokens[1]
	noteID := tokens[3]

	return &NoteResource{
		Project: project,
		NoteID:  noteID,
	}, nil
}

func (r *NoteResource) Name() string {
	return fmt.Sprintf("%s/notes/%s", r.Project, r.NoteID)
}
