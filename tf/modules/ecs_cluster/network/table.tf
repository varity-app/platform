resource "aws_route_table" "public" {
    vpc_id = aws_vpc.ecs.id

    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.ecs.id
    }

    tags = {
        Name = "Public Subnet Route Table"
    }
}

resource "aws_route_table_association" "public" {
    subnet_id = aws_subnet.ecs.id
    route_table_id = aws_route_table.public.id
}