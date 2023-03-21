# aactl

You can install `aactl` CLI using one of the following ways:

* [Go](#go)
* [Homebrew](#homebrew)
* [RHEL/CentOS](#rhelcentos)
* [Debian/Ubuntu](#debianubuntu)
* [Binary](#binary)

See the [release section](https://github.com/GoogleCloudPlatform/aactl/releases/latest) for `aactl` checksums and SBOMs.

### Go

If you have Go 1.17 or newer, you can install latest `aactl` using:

```shell
go install github.com/GoogleCloudPlatform/aactl/cmd/aactl@latest
```

### Homebrew

On Mac or Linux, you can install `aactl` with [Homebrew](https://brew.sh/):

```shell
brew tap GoogleCloudPlatform/aactl https://github.com/GoogleCloudPlatform/aactl
brew install aactl
```

New release will be automatically picked up when you run `brew upgrade`


### RHEL/CentOS

```shell
rpm -ivh https://github.com/GoogleCloudPlatform/aactl/releases/download/v$VERSION/aactl-$VERSION_Linux-amd64.rpm
```

### Debian/Ubuntu

```shell
wget https://github.com/aquasecurity/aactl/releases/download/v$VERSION/aactl-$VERSION_Linux-amd64.deb
sudo dpkg -i aactl-$VERSION_Linux-64bit.deb
```

### Binary 

You can also download the [latest release](https://github.com/GoogleCloudPlatform/aactl/releases/latest) version of `aactl` for your operating system/architecture from [here](https://github.com/GoogleCloudPlatform/aactl/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

> The official `aactl` releases include SBOMs

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
