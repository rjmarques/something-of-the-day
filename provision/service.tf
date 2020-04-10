resource "aws_ecs_service" "sotd-ecs-service" {
  	name            = "sotd-ecs-service"
  	cluster         = aws_ecs_cluster.sotd-ecs-cluster.id
    scheduling_strategy = "DAEMON"
  	task_definition = aws_ecs_task_definition.sotd.arn
}