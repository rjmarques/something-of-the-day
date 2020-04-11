resource "aws_autoscaling_group" "ecs-autoscaling-group" {
    name                        = "ecs-autoscaling-group"
    max_size                    = var.max_instance_size
    min_size                    = var.min_instance_size
    desired_capacity            = var.desired_capacity
    vpc_zone_identifier         = [aws_default_subnet.default.id]
    launch_configuration        = aws_launch_configuration.ecs-launch-configuration.name
    tag {
        key                     = "Name"
        value                   = "something-of-the-day"
        propagate_at_launch     = true
    }
}