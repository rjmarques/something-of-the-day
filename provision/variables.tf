## Main creds for AWS connection 
variable "ecs_cluster" {
  description = "ECS cluster name"
}

variable "account_id" {
  description = "AWS Accout ID"
}

variable "ecs_key_pair_name" {
  description = "EC2 instance key pair name"
}

variable "region" {
  description = "AWS region"
}

variable "availability_zone" {
  description = "AWS Subnet Availablity Zone"
}

## Autoscale Config

variable "max_instance_size" {
  description = "Maximum number of instances in the cluster"
}

variable "min_instance_size" {
  description = "Minimum number of instances in the cluster"
}

variable "desired_capacity" {
  description = "Desired number of instances in the cluster"
}

## Twitter account creds
variable "twitter_client_id" {
  description = "Twitter account client id used to connect to the Twitter API"
}

variable "twitter_client_secret" {
  description = "Twitter account client secret used to connect to the Twitter API"
}