resource "aws_secretsmanager_secret" "sotd_db_url" {
  name = "db-url"
}

resource "aws_secretsmanager_secret_version" "sotd_db_url_value" {
  secret_id     = aws_secretsmanager_secret.sotd_db_url.id
  secret_string = var.postgres_url
}

resource "aws_secretsmanager_secret" "sotd_twitter_client_id" {
  name = "twitter_client_id"
}

resource "aws_secretsmanager_secret_version" "sotd_twitter_client_id_value" {
  secret_id     = aws_secretsmanager_secret.sotd_twitter_client_id.id
  secret_string = var.twitter_client_id
}

resource "aws_secretsmanager_secret" "sotd_twitter_client_secret" {
  name = "twitter_client_secret"
}

resource "aws_secretsmanager_secret_version" "sotd_twitter_client_secret_value" {
  secret_id     = aws_secretsmanager_secret.sotd_twitter_client_secret.id
  secret_string = var.twitter_client_secret
}
