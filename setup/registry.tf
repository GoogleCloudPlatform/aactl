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

# Description: Creates a Google Artifact Registry for the project

# Artifact Registry
resource "google_artifact_registry_repository" "registry" {
  provider      = google-beta
  project       = var.project_id
  description   = "${var.name} artifacts registry"
  location      = var.location
  repository_id = var.name
  format        = "DOCKER"
}

# Role binding to allow publisher to publish images
resource "google_artifact_registry_repository_iam_member" "registry_role_binding" {
  provider   = google-beta
  project    = var.project_id
  location   = var.location
  repository = google_artifact_registry_repository.registry.name
  role       = "roles/artifactregistry.repoAdmin"
  member     = "serviceAccount:${google_service_account.github_actions_user.email}"
}

