terraform {
  required_providers {
    graalsystems = {
      source = "graalsystems/graalsystems"
      version = "1.0.5"
    }
  }
}

provider "graalsystems" {
  api_url = "http://172.24.240.1:4200/api/v1"
  auth_url = "http://172.24.240.1:8089"
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
  name = "my identity"
}

resource "graalsystems_job" "my_job" {
  project_id = graalsystems_project.my_project.id
  identity_id = graalsystems_identity.my_identity.id
  name = "my job"
  description = "a useful description job"
  labels = {
    "pipeline": "data-ingestion"
  }

  spark {
    instance_type = "Standard_General_G1_v1"
    main_class_name = "org.acme.MySparkJob"
  }
}