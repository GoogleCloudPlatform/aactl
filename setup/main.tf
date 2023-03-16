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

# Description: This file contains the Terraform code to enable the required GCP APIs for the project

# List of GCP APIs to enable in this project
locals {
  services = [
    "artifactregistry.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "container.googleapis.com",
    "containerregistry.googleapis.com",
    "containerscanning.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "servicecontrol.googleapis.com",
    "servicemanagement.googleapis.com",
  ]
}

# Data source to access GCP project metadata 
data "google_project" "project" {}


# Enable the required GCP APIs
resource "google_project_service" "default" {
  for_each = toset(local.services)

  project = var.project_id
  service = each.value

  disable_on_destroy = false
}