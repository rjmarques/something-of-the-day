resource "aws_ecr_repository" "sotd-repo" {
  name                 = "something-of-the-day"
  image_tag_mutability = "MUTABLE"
}
