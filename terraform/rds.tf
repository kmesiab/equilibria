resource "aws_db_instance" "mysql" {
  engine         = "mysql"
  engine_version = var.mysql_version

  identifier = var.app_name
  username   = aws_ssm_parameter.database_user.value
  password   = aws_ssm_parameter.database_password.value
  db_name    = aws_ssm_parameter.database_name.value

  allocated_storage      = 20
  storage_type           = "gp2"
  instance_class         = "db.t2.micro"  # Choose an appropriate instance size
  parameter_group_name   = "default.mysql8.0"
  db_subnet_group_name   = aws_db_subnet_group.mysql_subnet_group.name
  vpc_security_group_ids = [aws_security_group.mysql_sg.id]

  skip_final_snapshot = true
  publicly_accessible = true
}

# RDS subnet group places the RDS instance on both subnets
resource "aws_db_subnet_group" "mysql_subnet_group" {
  name       = "my-mysql-subnet-group"
  subnet_ids = [aws_subnet.receiver_subnet.id, aws_subnet.outbound_subnet.id]
}

# RDS security group
resource "aws_security_group" "mysql_sg" {
  name        = "mysql-security-group"
  description = "Allow inbound traffic on MySQL port"
  vpc_id      = aws_vpc.my_vpc.id

  # Inbound rule to allow traffic from the Lambda function
  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = [
      aws_subnet.receiver_subnet.cidr_block,
      aws_subnet.outbound_subnet.cidr_block  # Subnet CIDR blocks
    ]
  }

  #
  # My laptop
  #
  ingress {
    from_port        = 3306
    to_port          = 3306
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]  # Allow all IPv4 addresses
    ipv6_cidr_blocks = ["::/0"]  # Allow all IPv6 addresses
  }

  tags = {
    Name = "MySQLSG"
  }
}
