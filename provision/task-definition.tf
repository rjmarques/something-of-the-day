resource "aws_ecs_task_definition" "sotd" {
    family                = "sotd-definition"
    execution_role_arn    = aws_iam_role.ecs-task-execution-role.arn
    container_definitions = <<DEFINITION
[
  {
    "name": "something-of-the-day",
    "image": "661157442746.dkr.ecr.eu-west-2.amazonaws.com/something-of-the-day:latest",
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "memory": 128,
    "cpu": 1024,
    "secrets": [
      {
        "name": "POSTGRES_URL",
        "valueFrom": "${aws_secretsmanager_secret.sotd_db_url.arn}"
      },
      {
        "name": "CLIENT_ID",
        "valueFrom": "${aws_secretsmanager_secret.sotd_twitter_client_id.arn}"
      },
      {
        "name": "CLIENT_SECRET",
        "valueFrom": "${aws_secretsmanager_secret.sotd_twitter_client_secret.arn}"
      }
    ]
  }
]
DEFINITION
}