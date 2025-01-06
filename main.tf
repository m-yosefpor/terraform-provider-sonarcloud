terraform {
  required_providers {
    sonarcloud = {
      source  = "local/m-yosefpor/sonarcloud"
      version = "1.0.0"
    }
  }
}

provider "sonarcloud" {
  token        = var.sonarcloud_token
  organization = var.organization
}

resource "sonarcloud_project" "projects" {
  for_each = toset(var.projects)

  project_key               = "${var.organization}_${each.key}"
  name                      = each.key
  visibility                = "private"
  organization              = var.organization
  new_code_definition_type  = "days"
  new_code_definition_value = "30"
}

# resource "sonarcloud_project" "example2" {
#   project_key  = "bitvavo_matt_test2"
#   name         = "My Project 2"
#   visibility   = "private"
#   organization = "bitvavo"
#   new_code_definition_type = "days"
#   new_code_definition_value = "30"
# }


resource "sonarcloud_qualitygates_select" "projects" {
  for_each = sonarcloud_project.projects

  project_key     = each.value.project_key
  organization    = var.organization
  quality_gate_id = 117724
}

# resource "sonarcloud_qualitygates_select" "example2" {
#   project_key     = "bitvavo_frontend-hermes"
#   organization = var.organization
#   quality_gate_id = 117724
# }


variable "sonarcloud_token" {
}

variable "organization" {
  default = "bitvavo"
}


variable "projects" {
  type    = list(string)
  default = []
}