#
# Network LB
#  - LB -> Listener(s) -> Target Group(s) <- ASG
#
resource "aws_lb" "aslsp" {
  name                             = "aslsp"
  internal                         = false
  load_balancer_type               = "network"
  subnets                          = ["${data.aws_subnet.aslsp.*.id}"]
  enable_cross_zone_load_balancing = true
  enable_deletion_protection       = false

  tags {
    Name = "aslsp"
  }
}

# ASG uses target groups to connect with the network LB
# so we configure a group+listener for each port

# 80
#
resource "aws_lb_target_group" "aslsp-80" {
  name     = "aslsp-80"
  port     = 80
  protocol = "TCP"
  vpc_id   = "vpc-a27c82c6"
}

resource "aws_lb_listener" "aslsp-80" {
  load_balancer_arn = "${aws_lb.aslsp.arn}"
  port              = "80"
  protocol          = "TCP"

  default_action {
    target_group_arn = "${aws_lb_target_group.aslsp-80.arn}"
    type             = "forward"
  }
}

# 443
#
resource "aws_lb_target_group" "aslsp-443" {
  name     = "aslsp-443"
  port     = 443
  protocol = "TCP"
  vpc_id   = "vpc-a27c82c6"
}

resource "aws_lb_listener" "aslsp-443" {
  load_balancer_arn = "${aws_lb.aslsp.arn}"
  port              = "443"
  protocol          = "TCP"

  default_action {
    target_group_arn = "${aws_lb_target_group.aslsp-443.arn}"
    type             = "forward"
  }
}
