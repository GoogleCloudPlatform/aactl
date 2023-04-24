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

package vul

import (
	"context"
	"fmt"
	"sync"

	ca "cloud.google.com/go/containeranalysis/apiv1"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	"github.com/GoogleCloudPlatform/aactl/pkg/vul/convert"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	g "google.golang.org/genproto/googleapis/grafeas/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Import(ctx context.Context, options types.Options) error {
	opt, ok := options.(*types.VulnerabilityOptions)
	if !ok || opt == nil {
		return errors.New("valid options required")
	}
	if err := options.Validate(); err != nil {
		return errors.Wrap(err, "error validating options")
	}
	s, err := utils.NewFileSource(opt.File, opt.Source)
	if err != nil {
		return errors.Wrap(err, "error creating source")
	}

	converter, err := convert.GetConverter(opt.Format)
	if err != nil {
		return errors.Wrap(err, "error getting converter")
	}

	noteOccurrencesMap, err := converter(s)
	if err != nil {
		return errors.Wrap(err, "error converting source")
	}

	// TODO: Debug code
	//_ = deleteNoteOccurrences(ctx, opt, noteOccurrencesMap)

	log.Info().Msgf("found %d vulnerabilities", len(noteOccurrencesMap))

	return post(ctx, noteOccurrencesMap, opt)
}

func post(ctx context.Context, list types.NoteOccurrencesMap, opt *types.VulnerabilityOptions) error {
	if list == nil {
		return errors.New("expected non-nil result")
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for noteID, nocc := range list {
		wg.Add(1)
		go func(noteID string, nocc types.NoteOccurrences) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := postNoteOccurrences(ctx, opt.Project, noteID, nocc); err != nil {
				log.Error().Err(err).Msg("error posting notes")
				cancel()
			}
		}(noteID, nocc)
	}

	wg.Wait()

	return nil
}

// postNoteOccurrences creates new Notes and its associated Occurrences.
// Notes will be created only if it does not exist.
func postNoteOccurrences(ctx context.Context, projectID string, noteID string, nocc types.NoteOccurrences) error {
	if projectID == "" {
		return types.ErrMissingProject
	}

	// don't submit end-to-end test
	if projectID == types.TestProjectID {
		return nil
	}

	c, err := ca.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating client")
	}
	defer c.Close()

	p := fmt.Sprintf("projects/%s", projectID)

	// Create Note
	req := &g.CreateNoteRequest{
		Parent: p,
		NoteId: noteID,
		Note:   nocc.Note,
	}
	noteName := fmt.Sprintf("%s/notes/%s", p, noteID)
	_, err = c.GetGrafeasClient().CreateNote(ctx, req)
	if err != nil {
		// If note already exists, skip
		if status.Code(err) == codes.AlreadyExists {
			log.Debug().Msgf("already exists: %s", noteName)
		} else {
			return errors.Wrap(err, "error posting note")
		}
	}

	// Create Occurrences
	for _, o := range nocc.Occurrences {
		o.NoteName = noteName
		req := &g.CreateOccurrenceRequest{
			Parent:     p,
			Occurrence: o,
		}
		_, err := c.GetGrafeasClient().CreateOccurrence(ctx, req)
		if err != nil {
			// If occurrence already exists, skip
			if status.Code(err) == codes.AlreadyExists {
				log.Debug().Msgf("already exists: occurrence %s-%s",
					o.GetVulnerability().PackageIssue[0].AffectedPackage,
					o.GetVulnerability().PackageIssue[0].AffectedVersion.Name)
			} else {
				return errors.Wrap(err, "error posting occurrence")
			}
		}
	}

	return nil
}
