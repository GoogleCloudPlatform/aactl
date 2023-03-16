# Contributing

1. Sign one of the contributor license agreements below.
1. [Install Go](https://golang.org/doc/install).
1. Clone the repo:

   `git clone https://github.com/GoogleCloudPlatform/aactl.git`

1. Change into the checked out source:

   `cd aactl`

1. Fork the repo.
1. Set your fork as a remote:

   `git remote add fork https://github.com/GITHUB_USERNAME/golang-samples.git`

1. Make changes and commit to your fork. Initial commit messages should follow the
   [Conventional Commits](https://www.conventionalcommits.org/) style.
1. Send a pull request with your changes.
1. A maintainer will review the pull request and make comments. Prefer adding
   additional commits over amending and force-pushing since it can be difficult
   to follow code reviews when the commit history changes. 
   Commits will be squashed when they're merged.

## Testing

Before submitting PR, make sure the unit tests pass:

```shell
make test
```

And that there are no Go or YAML linting errors:

```shell
make lint
```

# Contributor License Agreements

Before we can accept your pull requests you'll need to sign a Contributor
License Agreement (CLA):

- **If you are an individual writing original source code** and **you own the
  intellectual property**, then you'll need to sign an [individual CLA][indvcla].
- **If you work for a company that wants to allow you to contribute your work**,
  then you'll need to sign a [corporate CLA][corpcla].

You can sign these electronically (just scroll to the bottom). After that,
we'll be able to accept your pull requests.

[gcloudcli]: https://developers.google.com/cloud/sdk/gcloud/
[indvcla]: https://developers.google.com/open-source/cla/individual
[corpcla]: https://developers.google.com/open-source/cla/corporate




