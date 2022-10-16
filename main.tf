provider "graalsystems" {
  api_url = "http://localhost:4200/api/v1"
  auth_url = "http://localhost:4200/api/v1"
  tenant = "XXX-XXX-XXX"
  auth_mode = "XXX-XXX-XXX"
  username = "XXX-XXX-XXX"
  password = "XXX-XXX-XXX"
  application_id = "XXX-XXX-XXX"
  application_secret = "XXX-XXX-XXX"
}

resource "project" "example_project" {
  name = "Example project"
  description = "This is an example project"
}