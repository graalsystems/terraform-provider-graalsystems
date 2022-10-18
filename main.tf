terraform {
  required_providers {
    graalsystems = {
      source = "graalsystems/graalsystems"
      version = "1.0.3"
    }
  }
}

provider "graalsystems" {
  api_url = "http://localhost:4200"
  auth_url = "http://localhost:8089"
  tenant = "platform-vincent-internal"
//  auth_mode = "credentials"
  username = "vdevillers"
  password = "devillerspwd"
//  application_id = "XXX-XXX-XXX"
//  application_secret = "XXX-XXX-XXX"
}

resource "graalsystems_project" "my_project" {
  name = "Example project"
  description = "This is an example project"
}


resource "graalsystems_identity" "my_identity" {
  name       = "my identity"
}

resource "graalsystems_job" "my_job" {
  project_id  = my_project.id
  identity_id  = my_identity.id
  name        = "my job"
  description = "a useful description job"
  tags        = ["spark", "just a tag"]

  spark {
    main_class_name = "org.acme.MySparkJob"
  }
}