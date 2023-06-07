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

package container

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	crname "github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/v2/pkg/cosign"
	"github.com/sigstore/cosign/v2/pkg/oci"
)

func getCosignOptions(ctx context.Context) (*cosign.CheckOpts, error) {
	rekorPubKeys, err := cosign.GetRekorPubs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", types.ErrInternal, err)
	}

	ctPubKeys, err := cosign.GetCTLogPubs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", types.ErrInternal, err)
	}

	roots, err := fulcio.GetRoots()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", types.ErrInternal, err)
	}
	intermediates, err := fulcio.GetIntermediates()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", types.ErrInternal, err)
	}

	return &cosign.CheckOpts{
		RootCerts:         roots,
		IntermediateCerts: intermediates,
		RekorPubKeys:      rekorPubKeys,
		CTLogPubKeys:      ctPubKeys,
	}, nil
}

func VerifyAndGetAttestations(ctx context.Context, resourceURI string) ([]oci.Signature, error) {
	// Get cosign.VerifyImageAttestations options
	opts, err := getCosignOptions(ctx)
	if err != nil {
		return nil, err
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
