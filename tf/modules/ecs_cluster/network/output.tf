output "vpc_id" {
  description = "ID of the ECS VPC"
  value       = aws_vpc.ecs.id
}

output "subnet_id" {
  description = "ID of ECS VPC main subnet"
  value       = aws_subnet.ecs.id
}

output "security_group" {
  description = "Name of ECS security group"
  value       = aws_security_group.ecs.id
}