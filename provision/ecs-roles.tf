###### ecs instance role ######
resource "aws_iam_role" "ecs-instance-role" {
    name                = "ecs-instance-role"
    path                = "/"
    assume_role_policy  = data.aws_iam_policy_document.ecs-instance-policy.json
}

data "aws_iam_policy_document" "ecs-instance-policy" {
    statement {
        actions = ["sts:AssumeRole"]

        principals {
            type        = "Service"
            identifiers = ["ec2.amazonaws.com"]
        }
    }
}

resource "aws_iam_role_policy_attachment" "ecs-instance-role-attachment" {
    role       = aws_iam_role.ecs-instance-role.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs-instance-profile" {
    name = "ecs-instance-profile"
    path = "/"
    role = aws_iam_role.ecs-instance-role.id
    provisioner "local-exec" {
      command = "sleep 10"
    }
}

###### ecs task execution role  ######
resource "aws_iam_role" "ecs-task-execution-role" {
    name                = "ecsTaskExecutionRole"
    path                = "/"
    assume_role_policy  = data.aws_iam_policy_document.ecs-task-execution-policy.json
}

data "aws_iam_policy_document" "ecs-task-execution-policy" {
    statement {
        actions = ["sts:AssumeRole"]

        principals {
            type        = "Service"
            identifiers = ["ecs-tasks.amazonaws.com"]
        }
    }
}

resource "aws_iam_role_policy_attachment" "ecs-task-execution-role-attachment" {
    role       = aws_iam_role.ecs-task-execution-role.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_policy" "ecs-secrets-access-policy" {
    name        = "SecretsAccessPolicy"
    description = "Allows ECS to access secrets defined in AWS Secret Manager"

    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetSecretValue"
            ],
            "Resource": [
                "${aws_secretsmanager_secret.sotd_db_url.arn}",
                "${aws_secretsmanager_secret.sotd_twitter_client_id.arn}",
                "${aws_secretsmanager_secret.sotd_twitter_client_secret.arn}"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-secrets-access-policy-attachment" {
    role       = aws_iam_role.ecs-task-execution-role.name
    policy_arn = aws_iam_policy.ecs-secrets-access-policy.arn
}
