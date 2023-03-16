package convert

import (
	"context"

	"github.com/GoogleCloudPlatform/aactl/pkg/convert/grype"
	"github.com/GoogleCloudPlatform/aactl/pkg/convert/snyk"
	"github.com/GoogleCloudPlatform/aactl/pkg/convert/trivy"
	"github.com/GoogleCloudPlatform/aactl/pkg/src"
	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/pkg/errors"
)

// VulnerabilityConverter is a function that converts a source to a list of AA notes.
type VulnerabilityConverter func(ctx context.Context, s *src.Source) (map[string]types.NoteOccurrences, error)

// GetConverter returns a vulnerability converter for the given format.
func GetConverter(format types.SourceFormat) (VulnerabilityConverter, error) {
	switch format {
	case types.SourceFormatSnykJSON:
		return snyk.Convert, nil
	case types.SourceFormatTrivyJSON:
		return trivy.Convert, nil
	case types.SourceFormatGrypeJSON:
		return grype.Convert, nil
	default:
		return nil, errors.Errorf("unimplemented conversion format: %s", format)
	}
}
