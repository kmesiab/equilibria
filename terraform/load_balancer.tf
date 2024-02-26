resource "aws_lb" "lb" {
  name               = "equilibria-lb"
  internal           = false
  load_balancer_type = "network"
  subnets            = [aws_subnet.receiver_subnet.id]

  enable_deletion_protection = false
}


# Define a target group
resource "aws_lb_target_group" "tg" {
  name     = "inbound-https-target-group"
  port     = 443
  protocol = "TCP"
  vpc_id   = aws_vpc.my_vpc.id

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_lb.lb.arn
  port              = 443
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.tg.arn
  }
}
