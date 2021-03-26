resource "aws_subnet" "ecs" {
    vpc_id = aws_vpc.ecs.id
    cidr_block = "10.0.0.0/24"

    tags = {
        Name = "ECS Public Subnet"
    }
}