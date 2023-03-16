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

# Description: Outputs for the deployment

output "REG_URI" {
  value       = "${google_artifact_registry_repository.registry.location}-docker.pkg.dev/${data.google_project.project.name}/${google_artifact_registry_repository.registry.name}"
  description = "Fully qualified Artifact Registry URI to use in Auth Actions."
}

output "IMG_NAME" {
  value       = google_artifact_registry_repository.registry.name
  description = "Image name to use in Auth Actions."
}

output "SA_EMAIL" {
  value       = google_service_account.github_actions_user.email
  description = "Service account to use in GitHub Actions."
}

output "PROVIDER_ID" {
  value       = google_iam_workload_identity_pool_provider.github_provider.name
  description = "Provider ID to use in Auth Actions."
}