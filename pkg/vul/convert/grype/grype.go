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

package grype

import (
	"fmt"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/utils"
	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
	g "google.golang.org/genproto/googleapis/grafeas/v1"
)

// Convert converts Snyk JSON to Grafeas Note/Occurrence format.
func Convert(s *utils.Source) (types.NoteOccurrencesMap, error) {
	if s == nil || s.Data == nil {
		return nil, errors.New("valid source required")
	}

	if !s.Data.Search("matches").Exists() {
		return nil, errors.New("unable to find vulnerabilities in source data")
	}

	list := make(types.NoteOccurrencesMap, 0)

	for _, v := range s.Data.Search("matches").Children() {
		// create note
		n := convertNote(s, v)

		// don't add notes with no CVSS score
		if n == nil || n.GetVulnerability().CvssScore == 0 {
			continue
		}
		noteID := utils.GetPrefixNoteName(n.GetShortDescription())

		// If cve is not found, add to map
		if _, ok := list[noteID]; !ok {
			list[noteID] = types.NoteOccurrences{Note: n}
		}
		nocc := list[noteID]
		occ := convertOccurrence(s, v, noteID)
		if occ != nil {
			nocc.Occurrences = append(nocc.Occurrences, occ)
		}
		list[noteID] = nocc
	}

	return list, nil
}

func convertOccurrence(s *utils.Source, v *gabs.Container, noteID string) *g.Occurrence {
	noteName := fmt.Sprintf("projects/%s/notes/%s", s.Project, noteID)

	// nvd vulnerability
	rvList := v.Search("relatedVulnerabilities").Children()
	var rv *gabs.Container
	for _, rvNode := range rvList {
		if rvNode.Search("namespace").Data().(string) == "nvd:cpe" {
			rv = rvNode
			break
		}
	}
	if rv == nil {
		return nil
	}
	cve := rv.Search("id").Data().(string)

	// cvssv2
	cvssList := rv.Search("cvss").Children()
	var cvss2, cvss3 *gabs.Container
	for _, cvss := range cvssList {
		switch cvss.Search("version").Data().(string) {
		case "2.0":
			cvss2 = cvss
		case "3.0", "3.1":
			cvss3 = cvss
		}
	}
	if cvss2 == nil {
		return nil
	}

	// Create Occurrence
	o := g.Occurrence{
		ResourceUri: fmt.Sprintf("https://%s", s.URI),
		NoteName:    noteName,
		Details: &g.Occurrence_Vulnerability{
			Vulnerability: &g.VulnerabilityOccurrence{
				ShortDescription: cve,
				LongDescription:  rv.Search("description").Data().(string),
				RelatedUrls: []*g.RelatedUrl{
					{
						Label: "Registry",
						Url:   s.URI,
					},
				},
				CvssVersion: g.CVSSVersion_CVSS_VERSION_2,
				CvssScore:   utils.ToFloat32(cvss2.Search("metrics", "baseScore").Data()),
				Severity:    utils.ToGrafeasSeverity(rv.Search("severity").Data().(string)),
				// TODO: What is the difference between severity and effective severity?
				EffectiveSeverity: utils.ToGrafeasSeverity(rv.Search("severity").Data().(string)),
			}},
	}

	// PackageIssues
	if len(v.Search("vulnerability", "fix", "versions").Children()) == 0 {
		o.GetVulnerability().PackageIssue = append(
			o.GetVulnerability().PackageIssue,
			getBasePackageIssue(v))
	} else {
		for _, version := range v.Search("vulnerability", "fix", "versions").Children() {
			pi := getBasePackageIssue(v)
			pi.FixedVersion = &g.Version{
				Name: version.Data().(string),
				Kind: g.Version_NORMAL,
			}
			o.GetVulnerability().PackageIssue = append(o.GetVulnerability().PackageIssue, pi)
		}
	}

	// CVSSv3
	if cvss3 != nil {
		o.GetVulnerability().Cvssv3 = utils.ToCVSS(
			utils.ToFloat32(cvss3.Search("metrics", "baseScore").Data()),
			cvss3.Search("vector").Data().(string),
		)
	}

	// References
	for _, r := range rv.Search("urls").Children() {
		o.GetVulnerability().RelatedUrls = append(o.GetVulnerability().RelatedUrls, &g.RelatedUrl{
			Url:   r.Data().(string),
			Label: "Url",
		})
	}
	return &o
}

func convertNote(s *utils.Source, v *gabs.Container) *g.Note {
	// nvd vulnerability
	rvList := v.Search("relatedVulnerabilities").Children()
	var rv *gabs.Container
	for _, rvNode := range rvList {
		if rvNode.Search("namespace").Data().(string) == "nvd:cpe" {
			rv = rvNode
			break
		}
	}
	if rv == nil {
		return nil
	}
	cve := rv.Search("id").Data().(string)

	// cvssv2
	cvssList := rv.Search("cvss").Children()
	var cvss2, cvss3 *gabs.Container
	for _, cvss := range cvssList {
		switch cvss.Search("version").Data().(string) {
		case "2.0":
			cvss2 = cvss
		case "3.0", "3.1":
			cvss3 = cvss
		}
	}
	if cvss2 == nil {
		return nil
	}

	// create note
	n := g.Note{
		ShortDescription: cve,
		LongDescription:  rv.Search("description").Data().(string),
		RelatedUrl: []*g.RelatedUrl{
			{
				Label: "Registry",
				Url:   s.URI,
			},
		},
		Type: &g.Note_Vulnerability{
			Vulnerability: &g.VulnerabilityNote{
				CvssVersion: g.CVSSVersion_CVSS_VERSION_2,
				CvssScore:   utils.ToFloat32(cvss2.Search("metrics", "baseScore").Data()),
				// Details in Notes are not populated since we will never see the full list
				Details: []*g.VulnerabilityNote_Detail{
					{
						AffectedCpeUri:  "N/A",
						AffectedPackage: "N/A",
					},
				},
				Severity: utils.ToGrafeasSeverity(rv.Search("severity").Data().(string)),
			},
		},
	} // end note

	// CVSSv3
	if cvss3 != nil {
		n.GetVulnerability().CvssV3 = utils.ToCVSSv3(
			utils.ToFloat32(cvss3.Search("metrics", "baseScore").Data()),
			cvss3.Search("vector").Data().(string),
		)
	}

	// References
	for _, r := range rv.Search("urls").Children() {
		n.RelatedUrl = append(n.RelatedUrl, &g.RelatedUrl{
			Url:   r.Data().(string),
			Label: "Url",
		})
	}

	return &n
}

func getBasePackageIssue(v *gabs.Container) *g.VulnerabilityOccurrence_PackageIssue {
	return &g.VulnerabilityOccurrence_PackageIssue{
		PackageType:     utils.ParsePackageType(v.Search("artifact", "language").Data().(string)),
		AffectedCpeUri:  v.Search("artifact", "cpes").Index(0).Data().(string),
		AffectedPackage: v.Search("artifact", "name").Data().(string),
		AffectedVersion: &g.Version{
			Name: v.Search("artifact", "version").Data().(string),
			Kind: g.Version_NORMAL,
		},
		FixedCpeUri:  v.Search("artifact", "cpes").Index(0).Data().(string),
		FixedPackage: v.Search("artifact", "name").Data().(string),
		FixedVersion: &g.Version{
			Kind: g.Version_MAXIMUM,
		},
	}
}
