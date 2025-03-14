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

package types

import (
	g "google.golang.org/genproto/googleapis/grafeas/v1"
)

// NoteOccurrences is a helper struct to hold Note and Occurrences.
type NoteOccurrences struct {
	// Note that the list of Occurrences points to.
	Note *g.Note

	// Occurrences that belong to the Note.
	Occurrences []*g.Occurrence
}

type NoteOccurrencesMap map[string]NoteOccurrences
