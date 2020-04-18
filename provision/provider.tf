provider "aws" {
  version = "~> 2.8"
  region = var.region
}

provider "heroku" {
  version = "~> 2.0"
}