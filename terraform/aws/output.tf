output "ecr_repository_url" {
  value = aws_ecr_repository.sotd-repo.repository_url
}

output "container_definition" {
  value = templatefile("${path.module}/container_definition.json", {
    repository_url        = aws_ecr_repository.sotd-repo.repository_url,
    db_url                = aws_ssm_parameter.sotd_db_url.arn,
    twitter_client_id     = aws_ssm_parameter.sotd_twitter_client_id.arn,
    twitter_client_secret = aws_ssm_parameter.sotd_twitter_client_secret.arn
  })
}

output "secrets_arns" {
  value = [
    aws_ssm_parameter.sotd_db_url.arn,
    aws_ssm_parameter.sotd_twitter_client_id.arn,
    aws_ssm_parameter.sotd_twitter_client_secret.arn
  ]
}
