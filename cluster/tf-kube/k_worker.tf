##
##	worker instances (aka nodes, minions)
##

resource "aws_instance" "worker" {
  count         = 3
  ami           = "${lookup(var.amis, var.region)}"
  instance_type = "${var.worker_inst_type}"

  subnet_id                   = "${aws_subnet.kube.id}"
  private_ip                  = "${cidrhost(var.vpc_cidr, 30 + count.index)}"
  associate_public_ip_address = true
  source_dest_check           = false

  availability_zone      = "${var.zone}"
  vpc_security_group_ids = ["${aws_security_group.kube.id}"]
  key_name               = "${var.default_keypair_name}"

  tags {
    Owner           = "${var.owner}"
    Name            = "worker-${count.index}"
    ansibleFilter   = "${var.ansibleFilter}"
    ansibleNodeType = "worker"
    ansibleNodeName = "worker${count.index}"
  }
}

##
##	Outputs
##

output "kube_workers" {
  value = "${join(",", aws_instance.worker.*.public_ip)}"
}
