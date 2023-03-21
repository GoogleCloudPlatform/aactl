# aactl

Google [Container Analysis (AA)](https://cloud.google.com/container-analysis/docs/container-analysis) service data import utility, supports OSS vulnerability scanner reports, SLSA provenance, and sigstore attestations.

> Installation instruction with support for Go, Homebrew, RHEL/CentOS, Debian/Ubuntu, and Binary are available [here](INSTALLATION.md).

## Usage 

`aactl` supports imports of two data types: `vulnerability` and `attestation`.

### Vulnerability 

To import vulnerabilities output by either [grype](https://github.com/anchore/grype), [snyk](https://github.com/snyk/cli), [trivy](https://github.com/aquasecurity/trivy) scanners, start by exporting the report in JSON format: 

* [grype](https://github.com/anchore/grype)

  `grype --add-cpes-if-none -s AllLayers -o json --file report.json $image`

* [snyk](https://github.com/snyk/cli)

  `snyk container test --app-vulns --json-file-output=report.json $image`

* [trivy](https://github.com/aquasecurity/trivy)

  `trivy image --format json --output report.json $image`

Once you have the vulnerability file, importing that file into AA using `aactl`:

```shell
aactl vulnerability --project $project \
                    --source $image \
                    --file report.json \
                    --format snyk
```

> The $image variable in the above example is the fully qualified URI of the image including its digest (e.g. `us-docker.pkg.dev/project/repo/image@sha256:397d453...`).

To review the imported vulnerabilities in GCP:

```shell
gcloud artifacts docker images list $repo \
  --show-occurrences \
  --format json \
  --occurrence-filter 'kind="VULNERABILITY" AND noteProjectId="$project" AND resource_url="$image" AND noteId="CVE-2005-2541"'
```

> You can also navigate to Artifact Registry to view the vulnerabilities there. 

### Attestation

In addition to vulnerabilities, `aactl` can also import [sigstore](https://github.com/sigstore) attestations:

```shell
aactl attestation --project $project \
                  --source $image \
```

> The $image variable in the above example is the fully qualified URI of the image including its digest (e.g. `us-docker.pkg.dev/project/repo/image@sha256:397d453...`).

## Contributing

Entirely new samples are not accepted. Bug fixes are welcome, either as pull
requests or as GitHub issues.

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute.

## Licensing

Code in this repository is licensed under the Apache 2.0. See [LICENSE](LICENSE).
