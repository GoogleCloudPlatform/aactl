# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Description: This file contains the resources required to federate a GitHub repository with a GCP project

# This is a list of roles that will be assigned to the GitHub federted user
locals {
  # List of roles that will be assigned to the GitHub federted user
  ci_roles = toset([
    "roles/artifactregistry.reader",
    "roles/artifactregistry.writer",
    "roles/containeranalysis.occurrences.editor",
    "roles/containeranalysis.notes.editor",
  ])
}

# Service account to be used for federated auth to publish to GCR (existing)
resource "google_service_account" "github_actions_user" {
  account_id   = "${var.name}-github-actions-user"
  display_name = "Service Account impersonated in ${var.git_repo} GitHub Actions"
}

# IAM policy bindings to the service account resources created by GitHub identify
resource "google_project_iam_member" "ci_role_bindings" {
  for_each = local.ci_roles
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.github_actions_user.email}"
}

# Identiy pool for GitHub action based identity's access to Google Cloud resources
resource "google_iam_workload_identity_pool" "github_pool" {
  provider                  = google-beta
  workload_identity_pool_id = "${var.name}-github-pool"
}

# Configuration for GitHub identiy provider
resource "google_iam_workload_identity_pool_provider" "github_provider" {
  provider                           = google-beta
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.aud"        = "assertion.aud"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }
  oidc {
    issuer_uri        = "https://token.actions.githubusercontent.com"
    allowed_audiences = []
  }
}

# IAM policy bindings to the service account resources created by GitHub identify
resource "google_service_account_iam_member" "pool_impersonation" {
  provider           = google-beta
  service_account_id = google_service_account.github_actions_user.id
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${var.git_repo}"
}
