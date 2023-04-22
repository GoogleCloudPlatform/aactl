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
	"google.golang.org/api/iterator"
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
	s, err := utils.NewFileSource(opt.Project, opt.File, opt.Source)
	if err != nil {
		return errors.Wrap(err, "error creating source")
	}

	c, err := convert.GetConverter(opt.Format)
	if err != nil {
		return errors.Wrap(err, "error getting converter")
	}

	list, err := c(s)
	if err != nil {
		return errors.Wrap(err, "error converting source")
	}

	// TODO: Debug code
	//_ = deleteNoteOccurrences(ctx, opt, list)

	log.Info().Msgf("found %d vulnerabilities", len(list))

	return post(ctx, list, opt)
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
				log.Error().Err(err).Msg("error posting notes & occurrences")
				cancel()
			}
		}(noteID, nocc)
	}

	wg.Wait()

	return nil
}

func post2(ctx context.Context, list types.NoteOccurrencesMap, opt *types.VulnerabilityOptions) error {
	if list == nil {
		return errors.New("expected non-nil result")
	}

	for noteID, nocc := range list {
		if err := postNoteOccurrences(ctx, opt.Project, noteID, nocc); err != nil {
			log.Error().Err(err).Msg("error posting notes & occurrences")

		}

	}

	return nil
}

// postNoteOccurrences creates new Notes and its associated Occurrences.
// Notes will be created only if it does not exist.
func postNoteOccurrences(ctx context.Context, projectID string, noteID string, nocc types.NoteOccurrences) error {
	if projectID == "" {
		return errors.New("projectID required")
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

	// Create/Update Occurrences
	for _, o := range nocc.Occurrences {
		if err := createOrUpdateOccurrence(ctx, p, noteID, o, c); err != nil {
			return errors.Wrap(err, "unable to create or update occurrence")
		}
		// TODO: PackageIssues should be merged
		break
	}

	return nil
}

func createOrUpdateOccurrence(ctx context.Context, p string, noteID string, o *g.Occurrence, c *ca.Client) error {
	// Create occurrence. If already exists, update.
	listReq := &g.ListOccurrencesRequest{
		Parent:   p,
		Filter:   fmt.Sprintf("noteId=\"%s\" AND resource_url=\"%s\"", noteID, o.GetResourceUri()),
		PageSize: 10,
	}

	var listRes []*g.Occurrence
	it := c.GetGrafeasClient().ListOccurrences(ctx, listReq)
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		listRes = append(listRes, resp)
	}

	switch len(listRes) {
	// If there were no occurrences, we create the occurrence.
	case 0:
		req := &g.CreateOccurrenceRequest{
			Parent:     p,
			Occurrence: o,
		}
		_, err := c.GetGrafeasClient().CreateOccurrence(ctx, req)
		if err != nil {
			return errors.Wrap(err, "error posting occurrence")
		}
	// If there was one occurrence, we update it.
	case 1:
		updateReq := &g.UpdateOccurrenceRequest{
			Name:       listRes[0].GetName(),
			Occurrence: o,
		}
		if _, err := c.GetGrafeasClient().UpdateOccurrence(ctx, updateReq); err != nil {
			return errors.Wrap(err, "error updating occurrence")
		}
	default:
		return errors.New("list occurrence expected to return one " +
			"occurrence but more than one was returned")
	}
	return nil
}
