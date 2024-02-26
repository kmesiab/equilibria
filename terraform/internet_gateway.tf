#
#
#
resource "aws_internet_gateway" "equilibria_internet_gateway" {
  vpc_id = aws_vpc.my_vpc.id

  tags = {
    Name = "EquilibriaInternetGateway"
  }
}

#
# Routes all outbound traffic to the internet gateway
#
resource "aws_route_table" "outbound_route_table" {
  vpc_id = aws_vpc.my_vpc.id

  # Open to the internet by way of the internet gateway
  # But no inbound routes back to the lambda
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.equilibria_internet_gateway.id
  }

  tags = {
    Name = "Outbound"
    Subnet = "outbound_subnet"
    Type = "Public"
    Description = "Routes outbound traffic to the internet gateway, from the send sms lambda"
  }
}

#
# Subnet 2 has access to the internet, so we give it a route table
# that points to the internet gateway  It only has outbound traffic
#
resource "aws_route_table_association" "outbound_subnet_association" {
  subnet_id      = aws_subnet.outbound_subnet.id
  route_table_id = aws_route_table.outbound_route_table.id
}
