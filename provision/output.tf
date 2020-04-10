output "heroku_db_url" {
  value = data.heroku_app.sotd.config_vars.DATABASE_URL
}

output "ecr_repository_url" {
  value = data.aws_ecr_repository.sotd-repo.repository_url
}