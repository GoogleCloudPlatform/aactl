package container

import (
	"context"

	crname "github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/oci"
)

func VerifyAndGetAttestations(ctx context.Context, resourceURI string) ([]oci.Signature, error) {
	// Attestations that chain upto Fulcio's root cert
	// TODO: Should handle case where attestations are not signed using Fulcio's
	//   short-lived keys. ie. Don't verify and return attestations.
	roots, err := fulcio.GetRoots()
	if err != nil {
		return nil, err
	}
	opts := &cosign.CheckOpts{
		RootCerts: roots,
	}

	// Verify and return attestations
	resourceRef, err := crname.ParseReference(resourceURI)
	if err != nil {
		return nil, err
	}

	atts, _, err := cosign.VerifyImageAttestations(ctx, resourceRef, opts)
	if err != nil {
		return nil, err
	}

	return atts, err
}
