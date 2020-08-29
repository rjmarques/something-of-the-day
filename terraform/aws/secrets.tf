resource "aws_ssm_parameter" "sotd_db_url" {
  name        = "/production/sotd/database/db_url"
  description = "PG connection URL"
  type        = "SecureString"
  value       = var.postgres_url

  tags = {
    environment = "production"
  }
}

resource "aws_ssm_parameter" "sotd_twitter_client_id" {
  name        = "/production/sotd/database/twitter_client_id"
  description = "Twitter Client ID"
  type        = "String"
  value       = var.twitter_client_id

  tags = {
    environment = "production"
  }
}

resource "aws_ssm_parameter" "sotd_twitter_client_secret" {
  name        = "/production/sotd/database/twitter_client_secret"
  description = "Twitter Client Secret"
  type        = "SecureString"
  value       = var.twitter_client_secret

  tags = {
    environment = "production"
  }
}
