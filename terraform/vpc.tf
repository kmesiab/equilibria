resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"

  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "EquilibriaVPC"
  }
}

#
# Subnet 1
#
# This network only allows inbound traffic.
#
resource "aws_subnet" "receiver_subnet" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.10.0/24"  # Adjusted non-overlapping CIDR block
  availability_zone = "us-west-2a"

  tags = {
    Name = "receiver-subnet"
    Description = "Inbound traffic from the API Gateway to the receiver sms and user management lambdas"
  }
}

#
# Subnet 2
#
resource "aws_subnet" "outbound_subnet" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.30.0/24"  # Adjusted non-overlapping CIDR block
  availability_zone = "us-west-2b"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "sender-subnet"
    Description = "Outbound traffic from the send sms lambda to the internet"
  }
}

resource "aws_vpc_endpoint" "ssm_vpc_endpoint" {
  vpc_id            = aws_vpc.my_vpc.id
  service_name      = "com.amazonaws.${var.region}.ssm"
  vpc_endpoint_type = "Interface"

  subnet_ids = [aws_subnet.receiver_subnet.id, aws_subnet.outbound_subnet.id]

  security_group_ids = [
    aws_security_group.login_lambda_sg.id
  ]

  private_dns_enabled = true
}

