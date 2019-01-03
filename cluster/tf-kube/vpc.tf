##
##	VPC
##

resource "aws_vpc" "kube" {
  cidr_block           = "${var.vpc_cidr}"
  enable_dns_hostnames = true

  tags {
    Name  = "${var.vpc_name}"
    Owner = "${var.owner}"
  }
}

# DHCP Options are not actually required, being identical to the Default Option Set
resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name         = "${region}.compute.internal"
  domain_name_servers = ["AmazonProvidedDNS"]

  tags {
    Name  = "${var.vpc_name}"
    Owner = "${var.owner}"
  }
}

resource "aws_vpc_dhcp_options_association" "dns_resolver" {
  vpc_id          = "${aws_vpc.kube.id}"
  dhcp_options_id = "${aws_vpc_dhcp_options.dns_resolver.id}"
}

##
##	SSH
##

resource "aws_key_pair" "maroda" {
  key_name   = "maroda"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDtHDQHePgT3SD/5y+xCFKbMDBjUKFrzADTKPYjel6GHpPe58OeORJZWtC6hPKVae4Vjn1CHut3TCymaacdGDtJ3LKZdCjYn93VGPoRNttyODMNo08mY3gt0qyt65hgRp3JGvPHGlscedqVTDLaqOSuslkUXQhVoIHCtgUQwzIG6ADoqPGFJOgGEbPWYY95MObb9uNpVkOmZ+T6+fuee+yuSB4dUq8SnrqAoKuamhD+FMny3zF3C+aWx6Z7fXKmD0/X3+BllpdEnKOsLdQXJvvE69GQxJoIsEgl/1ZICILV8QxFFxtNN/RFws87ouz7LurJdUa4VBNXeQo6ymEGizxx matt@oscillator.localdomain"
}

##
##	Subnets
##

# Subnet (public)
resource "aws_subnet" "kube" {
  vpc_id            = "${aws_vpc.kube.id}"
  cidr_block        = "${var.vpc_cidr}"
  availability_zone = "${var.zone}"

  tags {
    Name  = "kube"
    Owner = "${var.owner}"
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.kube.id}"

  tags {
    Name  = "kube"
    Owner = "${var.owner}"
  }
}

##
##	Routing
##

resource "aws_route_table" "kube" {
  vpc_id = "${aws_vpc.kube.id}"

  # Default route thru internet gw
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.gw.id}"
  }

  tags {
    Name  = "kube"
    Owner = "${var.owner}"
  }
}

resource "aws_route_table_association" "kube" {
  subnet_id      = "${aws_subnet.kube.id}"
  route_table_id = "${aws_route_table.kube.id}"
}

##
##	Sec
##

resource "aws_security_group" "kube" {
  vpc_id = "${aws_vpc.kube.id}"
  name   = "kube"

  # Allow all outbound
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow ICMP from control host IP
  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["${var.control_cidr}"]
  }

  # Allow all internal
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.vpc_cidr}"]
  }

  # Allow all traffic from the API ELB
  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = ["${aws_security_group.kube_api.id}"]
  }

  # Allow all traffic from control host IP
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.control_cidr}"]
  }

  tags {
    Name  = "kube"
    Owner = "${var.owner}"
  }
}
