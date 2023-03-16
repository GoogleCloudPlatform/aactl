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

# Description: This file contains the provider configuration for the project

# Configure the Google Cloud provider
terraform {
  required_version = ">= 1.1"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.4"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 4.4"
    }
  }
}

# Configure the Google Cloud provider
provider "google" {
  project = var.project_id
}

# Configure the beta version of Google Cloud provider
provider "google-beta" {
  project = var.project_id
}