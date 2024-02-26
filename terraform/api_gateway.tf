 #
# The API Gateway
#
resource "aws_api_gateway_rest_api" "api_gateway" {
  name        = "EquilibriaAPI"
  description = "API Gateway exposes routes to enable a twilio callback"
}

resource "aws_api_gateway_method_settings" "settings" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  stage_name  = aws_api_gateway_stage.api_stage.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled    = true
    logging_level      = "INFO"
    data_trace_enabled = true
  }
}

resource "aws_route_table" "receiver_route_table" {
  vpc_id = aws_vpc.my_vpc.id

  # Local route for VPC traffic
  route {
    cidr_block = aws_vpc.my_vpc.cidr_block
    gateway_id = "local"
  }

  tags = {
    Name = "receiver_route_table"
  }
}

resource "aws_route_table_association" "receiver_subnet_association" {
  subnet_id      = aws_subnet.receiver_subnet.id
  route_table_id = aws_route_table.receiver_route_table.id
}

resource "aws_api_gateway_vpc_link" "my_vpc_link" {
  name        = "EquilibriaVPCLink"
  target_arns = [aws_lb.lb.arn]  # ARN of the Network Load Balancer
  description = "Links to the load balancer"
}
