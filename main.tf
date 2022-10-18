terraform {
  required_providers {
    graalsystems = {
      source = "graalsystems/graalsystems"
    }
  }
  required_version = ">= 0.13"
}

provider "graalsystems" {
  api_url = "http://localhost:4200"
  auth_url = "http://localhost:8089"
  tenant = "platform-vincent-internal"
  auth_mode = "credentials"
  username = "vdevillers"
  password = "devillerspwd"
//  application_id = "XXX-XXX-XXX"
//  application_secret = "XXX-XXX-XXX"
}

resource "project" "example_project" {
  name = "Example project"
  description = "This is an example project"
}