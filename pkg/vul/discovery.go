package vul

import (
	"context"
	"fmt"
	"time"

	ca "cloud.google.com/go/containeranalysis/apiv1"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	g "google.golang.org/genproto/googleapis/grafeas/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func updateDiscoveryNoteAndOcc(ctx context.Context, projectID string, resourceURL string) error {
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
	discoveryNoteID := "aactl-PACKAGE_VULNERABILITY"
	if err := updateDiscoveryNote(ctx, p, discoveryNoteID, c); err != nil {
		return err
	}

	return updateDiscoveryOcc(ctx, p, discoveryNoteID, resourceURL, c)
}

func updateDiscoveryNote(ctx context.Context, parent string, discoveryNoteID string, c *ca.Client) error {
	// Create Note
	req := &g.CreateNoteRequest{
		Parent: parent,
		NoteId: discoveryNoteID,
		Note: &g.Note{
			ShortDescription: "aactl discovery note",
			Kind:             g.NoteKind_DISCOVERY,
			Type: &g.Note_Discovery{
				Discovery: &g.DiscoveryNote{AnalysisKind: g.NoteKind_DISCOVERY},
			},
		},
	}
	noteName := fmt.Sprintf("%s/notes/%s", parent, discoveryNoteID)
	_, err := c.GetGrafeasClient().CreateNote(ctx, req)
	if err != nil {
		// If note already exists, skip
		if status.Code(err) == codes.AlreadyExists {
			log.Debug().Msgf("already exists: %s", noteName)
		} else {
			return errors.Wrap(err, "error creating discovery note")
		}
	}
	return nil
}

func updateDiscoveryOcc(ctx context.Context, parent string, discoveryNoteID string, resourceURL string, c *ca.Client) error {
	noteName := fmt.Sprintf("%s/notes/%s", parent, discoveryNoteID)
	occ := &g.Occurrence{
		Kind:        g.NoteKind_DISCOVERY,
		ResourceUri: resourceURL,
		NoteName:    noteName,
		Details: &g.Occurrence_Discovery{
			Discovery: &g.DiscoveryOccurrence{
				ContinuousAnalysis: g.DiscoveryOccurrence_INACTIVE,
				AnalysisStatus:     g.DiscoveryOccurrence_COMPLETE,
				LastScanTime:       timestamppb.New(time.Now()),
			},
		},
	}

	listOccReq := &g.ListOccurrencesRequest{
		Parent: parent,
		Filter: fmt.Sprintf("resourceUrl=\"%s\" AND kind=\"DISCOVERY\" AND noteId=\"%s\"", resourceURL, discoveryNoteID),
	}

	var listRes []*g.Occurrence
	it := c.GetGrafeasClient().ListOccurrences(ctx, listOccReq)
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
			Parent:     parent,
			Occurrence: occ,
		}
		_, err := c.GetGrafeasClient().CreateOccurrence(ctx, req)
		if err != nil {
			return errors.Wrap(err, "error posting discovery occurrence")
		}
	// If there was one occurrence, we update it.
	case 1:
		updateReq := &g.UpdateOccurrenceRequest{
			Name:       listRes[0].GetName(),
			Occurrence: occ,
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
