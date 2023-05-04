# aactl

If you have Go 1.17 or newer, you can install latest `aactl` using:

```shell
go install github.com/GoogleCloudPlatform/aactl/cmd/aactl@latest
```

## Prerequisites 

Since you are interested in `aactl`, you probably already have GCP account and project. If not, you learn about creating and managing projects [here](https://cloud.google.com/resource-manager/docs/creating-managing-projects). The other prerequisites include:

### APIs

`aactl` also depends on a few GCP service APIs. To enable these, run:

```shell
gcloud services enable containeranalysis.googleapis.com
```

### Roles

Make sure you have the following Identity and Access Management (IAM) roles in each project: 

> Learn how to grant multiple IAM roles to a user [here](https://cloud.google.com/iam/docs/granting-changing-revoking-access#multiple-roles)

```shell
roles/artifactregistry.reader
roles/containeranalysis.occurrences.editor
roles/containeranalysis.notes.editor
```

If you experience any issues, you can see the project level policy using following command:

```shell
gcloud projects get-iam-policy $PROJECT_ID --format=json > policy.json
```

### Credentials

When running locally, `aactl` will look for Google account credentials in one of the well-known locations. To ensure your Application Default Credentials (ADC) are used by the `aactl` run this `gcloud` command and follow the prompts:

```shell
gcloud auth application-default login
```

> More about ADC [here](https://cloud.google.com/docs/authentication/provide-credentials-adc)


## Licensing

Code in this repository is licensed under the Apache 2.0. See [LICENSE](LICENSE).
