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

package provenance02

import (
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/aactl/pkg/provenance"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/pkg/errors"
	dsselib "github.com/secure-systems-lab/go-securesystemslib/dsse"
	g "google.golang.org/genproto/googleapis/grafeas/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// Convert converts a provenance statement to a Grafeas note and occurrence.
func Convert(nr utils.NoteResource, resourceURL string, env *provenance.Envelope) (*g.Note, *g.Occurrence, error) {
	prov, err := getProvenance(env)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error getting provenance")
	}

	n := g.Note{
		Name: nr.NoteID,
		Type: &g.Note_Build{
			Build: &g.BuildNote{
				BuilderVersion: prov.Predicate.Builder.ID,
			},
		},
	} // end note

	predicate, err := getPredicate(prov)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error getting provenance predicate")
	}

	o := g.Occurrence{
		ResourceUri: fmt.Sprintf("https://%s", resourceURL),
		NoteName:    nr.Name(),
		Details: &g.Occurrence_Build{
			Build: &g.BuildOccurrence{
				IntotoStatement: &g.InTotoStatement{
					Type:          prov.Type,
					PredicateType: prov.PredicateType,
					Subject:       getSubject(prov.Subject),
					Predicate:     predicate,
				},
			},
		},
		Envelope: &g.Envelope{
			Payload:     []byte(env.Payload),
			PayloadType: env.PayloadType,
			Signatures:  getSignatures(env.Envelope),
		},
	}
	return &n, &o, nil
}

func getPredicate(prov *intoto.ProvenanceStatementSLSA02) (*g.InTotoStatement_SlsaProvenanceZeroTwo, error) {
	pred := prov.Predicate

	inv, err := getPredicateInvocation(prov)
	if err != nil {
		return nil, err
	}

	p := g.InTotoStatement_SlsaProvenanceZeroTwo{
		SlsaProvenanceZeroTwo: &g.SlsaProvenanceZeroTwo{
			Builder: &g.SlsaProvenanceZeroTwo_SlsaBuilder{
				Id: pred.Builder.ID,
			},
			BuildType:  pred.BuildType,
			Invocation: inv,
			Metadata:   getPredicateMetadata(prov),
			Materials:  getPredicateMaterials(prov),
		},
	}

	return &p, nil
}

func getPredicateInvocation(prov *intoto.ProvenanceStatementSLSA02) (*g.SlsaProvenanceZeroTwo_SlsaInvocation, error) {
	inv := prov.Predicate.Invocation

	// Parameters
	pj, _ := json.Marshal(inv.Parameters)
	parameters := structpb.Struct{}
	err := parameters.UnmarshalJSON(pj)
	if err != nil {
		return nil, err
	}

	// Environment
	ej, _ := json.Marshal(inv.Environment)
	environment := structpb.Struct{}
	err = environment.UnmarshalJSON(ej)
	if err != nil {
		return nil, err
	}

	i := g.SlsaProvenanceZeroTwo_SlsaInvocation{
		ConfigSource: &g.SlsaProvenanceZeroTwo_SlsaConfigSource{
			Uri:        inv.ConfigSource.URI,
			Digest:     inv.ConfigSource.Digest,
			EntryPoint: inv.ConfigSource.EntryPoint,
		},
		Parameters:  &parameters,
		Environment: &environment,
	}

	return &i, nil
}

func getPredicateMetadata(prov *intoto.ProvenanceStatementSLSA02) *g.SlsaProvenanceZeroTwo_SlsaMetadata {
	m := prov.Predicate.Metadata
	metadata := g.SlsaProvenanceZeroTwo_SlsaMetadata{
		BuildInvocationId: m.BuildInvocationID,
		BuildStartedOn:    utils.ToGRPCTime(m.BuildStartedOn),
		BuildFinishedOn:   utils.ToGRPCTime(m.BuildFinishedOn),
		Completeness: &g.SlsaProvenanceZeroTwo_SlsaCompleteness{
			Parameters:  m.Completeness.Parameters,
			Environment: m.Completeness.Environment,
			Materials:   m.Completeness.Materials,
		},
		Reproducible: true,
	}

	return &metadata
}

func getPredicateMaterials(prov *intoto.ProvenanceStatementSLSA02) []*g.SlsaProvenanceZeroTwo_SlsaMaterial {
	pred := prov.Predicate
	materials := []*g.SlsaProvenanceZeroTwo_SlsaMaterial{}
	for _, m := range pred.Materials {
		materials = append(materials, &g.SlsaProvenanceZeroTwo_SlsaMaterial{
			Uri:    m.URI,
			Digest: m.Digest,
		})
	}
	return materials
}

func getSubject(subjects []intoto.Subject) []*g.Subject {
	s := []*g.Subject{}
	for _, subject := range subjects {
		s = append(s, &g.Subject{
			Name:   subject.Name,
			Digest: subject.Digest,
		})
	}
	return s
}

func getSignatures(env *dsselib.Envelope) []*g.EnvelopeSignature {
	signatures := []*g.EnvelopeSignature{}
	for _, s := range env.Signatures {
		signatures = append(signatures, &g.EnvelopeSignature{
			Sig:   []byte(s.Sig),
			Keyid: s.KeyID,
		})
	}
	return signatures
}

func getProvenance(env *provenance.Envelope) (*intoto.ProvenanceStatementSLSA02, error) {
	prov := intoto.ProvenanceStatementSLSA02{}
	if err := json.Unmarshal(env.DecodedPayload, &prov); err != nil {
		return nil, errors.Wrap(err, "error decoding DSSE env")
	}

	return &prov, nil
}
