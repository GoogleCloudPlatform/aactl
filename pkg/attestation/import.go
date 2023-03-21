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

package attestation

import (
	"context"
	"crypto/sha256"
	"fmt"

	ca "cloud.google.com/go/containeranalysis/apiv1"
	"github.com/GoogleCloudPlatform/aactl/pkg/container"
	"github.com/GoogleCloudPlatform/aactl/pkg/provenance"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	g "google.golang.org/genproto/googleapis/grafeas/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Import imports attestation metadata from a source.
func Import(ctx context.Context, opt *types.AttestationOptions) error {
	if opt == nil {
		return errors.New("options required")
	}
	if err := opt.Validate(); err != nil {
		return errors.Wrap(err, "error validating options")
	}

	resourceURL, err := container.GetFullURL(opt.Source)
	if err != nil {
		return errors.Wrap(err, "error getting full url")
	}
	log.Info().Msgf("Resource URL: %s", resourceURL)

	nr := utils.NoteResource{
		Project: fmt.Sprintf("projects/%s", opt.Project),
		NoteID:  fmt.Sprintf("intoto_%x", sha256.Sum256([]byte(resourceURL))),
	}

	envs, err := provenance.GetVerifiedEnvelopes(ctx, resourceURL)
	if err != nil {
		return errors.Wrap(err, "error unpacking message")
	}

	//_ = deleteNoteOccurrences(ctx, nr, resourceURL)

	err = importEnvelopes(ctx, envs, nr, resourceURL)
	if err != nil {
		return errors.Wrap(err, "error importing envelopes")
	}

	return nil
}

func importEnvelopes(ctx context.Context, envs []*provenance.Envelope, nr utils.NoteResource, resourceURL string) error {
	c, err := ca.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating client")
	}
	defer c.Close()

	for _, env := range envs {
		n, o, err := Convert(nr, resourceURL, env)
		if err != nil {
			return errors.Wrap(err, "error importing envelopes")
		}

		err = postNote(ctx, c, nr, n)
		if err != nil {
			return errors.Wrap(err, "error posting Note")
		}

		err = postOccurrence(ctx, c, nr, o)
		if err != nil {
			return errors.Wrap(err, "error posting Occurrence")
		}
	}

	return nil
}

func postNote(ctx context.Context, c *ca.Client, nr utils.NoteResource, n *g.Note) error {
	// Create Note
	req := &g.CreateNoteRequest{
		Parent: nr.Project,
		NoteId: nr.NoteID,
		Note:   n,
	}
	_, err := c.GetGrafeasClient().CreateNote(ctx, req)
	if err != nil {
		// If note already exists, skip
		if status.Code(err) == codes.AlreadyExists {
			log.Info().Msgf("Already Exists: %s", nr.Name())
		} else {
			return errors.Wrap(err, "error posting note")
		}
	} else {
		log.Info().Msgf("Created Note: %s", nr.Name())
	}

	return nil
}

func postOccurrence(ctx context.Context, c *ca.Client, nr utils.NoteResource, o *g.Occurrence) error {
	// Create Occurrence
	oreq := &g.CreateOccurrenceRequest{
		Parent:     nr.Project,
		Occurrence: o,
	}
	occ, err := c.GetGrafeasClient().CreateOccurrence(ctx, oreq)
	if err != nil {
		// If occurrence already exists, skip
		if status.Code(err) == codes.AlreadyExists {
			log.Info().Msgf("Already Exists: Occurrence")
		} else {
			return errors.Wrap(err, "error posting occurrence")
		}
	} else {
		log.Info().Msgf("Created Occurrence: %s", occ.Name)
	}

	return nil
}

// deleteNoteOccurrences deletes notes and occurrences. Used for debugging.
// nolint:unused
func deleteNoteOccurrences(ctx context.Context, nr utils.NoteResource, resourceURL string) error {
	c, err := ca.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating client")
	}
	defer c.Close()

	// Delete Notes

	dr := &g.DeleteNoteRequest{
		Name: nr.Name(),
	}
	_ = c.GetGrafeasClient().DeleteNote(ctx, dr)

	// Delete Occurrences
	req := &g.ListOccurrencesRequest{
		Parent:   nr.Project,
		Filter:   fmt.Sprintf("resource_url=\"https://%s\"", resourceURL),
		PageSize: 1000,
	}
	it := c.GetGrafeasClient().ListOccurrences(ctx, req)
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		dr := &g.DeleteOccurrenceRequest{
			Name: resp.Name,
		}
		_ = c.GetGrafeasClient().DeleteOccurrence(ctx, dr)
	}

	return nil
}
