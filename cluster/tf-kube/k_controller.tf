##
##	control plane instances
##

resource "aws_instance" "controller" {
  count         = 3
  ami           = "${lookup(var.amis, var.region)}"
  instance_type = "${var.controller_instance_type}"

  iam_instance_profile = "${aws_iam_instance_profile.kube.id}"

  subnet_id                   = "${aws_subnet.kube.id}"
  private_ip                  = "${cidrhost(var.vpc_cidr, 20 + count.index)}"
  associate_public_ip_address = true
  source_dest_check           = false

  availability_zone      = "${var.zone}"
  vpc_security_group_ids = ["${aws_security_group.kube.id}"]
  key_name               = "${var.default_keypair_name}"

  tags {
    Owner           = "${var.owner}"
    Name            = "controller-${count.index}"
    ansibleFilter   = "${var.ansibleFilter}"
    ansibleNodeType = "controller"
    ansibleNodeName = "controller${count.index}"
  }
}

##
##	API LB
##

resource "aws_elb" "kube_api" {
  name                      = "${var.elb_name}"
  instances                 = ["${aws_instance.controller.*.id}"]
  subnets                   = ["${aws_subnet.kube.id}"]
  cross_zone_load_balancing = false

  security_groups = ["${aws_security_group.kube_api.id}"]

  listener {
    lb_port           = 6443
    instance_port     = 6443
    lb_protocol       = "TCP"
    instance_protocol = "TCP"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 15
    target              = "HTTP:8080/healthz"
    interval            = 30
  }

  tags {
    Name  = "kube"
    Owner = "${var.owner}"
  }
}

##
##	Sec
##

resource "aws_security_group" "kube_api" {
  vpc_id = "${aws_vpc.kube.id}"
  name   = "kube-api"

  # Allow inbound traffic to the port used by Kubernetes API HTTPS
  ingress {
    from_port   = 6443
    to_port     = 6443
    protocol    = "TCP"
    cidr_blocks = ["${var.control_cidr}"]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Owner = "${var.owner}"
    Name  = "kube-api"
  }
}

##
##	Outputs
##

output "kube_api_dns_name" {
  value = "${aws_elb.kube_api.dns_name}"
}
