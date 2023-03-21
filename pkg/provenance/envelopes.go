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

package provenance

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/GoogleCloudPlatform/aactl/pkg/container"
	"github.com/GoogleCloudPlatform/aactl/pkg/dsse"
)

type Envelope struct {
	*dsse.DecodedEnvelope
	IntotoType          string
	IntotoPredicateType string
}

func GetVerifiedEnvelopes(ctx context.Context, resourceURI string) ([]*Envelope, error) {
	atts, err := container.VerifyAndGetAttestations(ctx, resourceURI)
	if err != nil {
		return nil, errors.Wrap(err, "error getting verified envelopes")
	}

	envs := []*Envelope{}
	for _, att := range atts {
		env, err := dsse.AttestationToEnvelope(att)
		if err != nil {
			return nil, errors.Wrap(err, "error getting verified envelopes")
		}

		decodedEnv, err := dsse.GetDecodedEnvelope(env)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding verified envelopes")
		}

		// Check in-toto version and slsa predicate type
		penv, err := getEnvelope(decodedEnv)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding verified envelopes")
		}

		log.Debug().Msgf("In-Toto Type (%s), PredicateType (%s)", penv.IntotoType, penv.IntotoPredicateType)

		// TODO: Currently only one slsa version supported
		if penv.IntotoType != "https://in-toto.io/Statement/v0.1" || penv.IntotoPredicateType != "https://slsa.dev/provenance/v0.2" {
			continue
		}

		envs = append(envs, penv)
	}

	return envs, nil
}

func getEnvelope(env *dsse.DecodedEnvelope) (*Envelope, error) {
	pred := struct {
		Type          string `json:"_type"`
		PredicateType string `json:"predicateType"`
	}{}
	if err := json.Unmarshal(env.DecodedPayload, &pred); err != nil {
		return nil, errors.Wrap(err, "error decoding DSSE env")
	}

	penv := Envelope{
		DecodedEnvelope:     env,
		IntotoType:          pred.Type,
		IntotoPredicateType: pred.PredicateType,
	}
	return &penv, nil
}
