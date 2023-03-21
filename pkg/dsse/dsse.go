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

package dsse

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
	dsselib "github.com/secure-systems-lab/go-securesystemslib/dsse"
	"github.com/sigstore/cosign/pkg/oci"
)

type DecodedEnvelope struct {
	*dsselib.Envelope
	DecodedPayload []byte
}

func AttestationToEnvelope(att oci.Signature) (*dsselib.Envelope, error) {
	payload, err := att.Payload()
	if err != nil {
		return nil, errors.Wrap(err, "error getting payload from attestation")
	}

	env := &dsselib.Envelope{}
	if err = json.Unmarshal(payload, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling DSSE env")
	}

	return env, nil
}

func GetDecodedEnvelope(env *dsselib.Envelope) (*DecodedEnvelope, error) {
	pyld, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding DSSE env")
	}

	de := &DecodedEnvelope{
		Envelope:       env,
		DecodedPayload: pyld,
	}

	return de, nil
}
