#
# AZs
#
data "aws_availability_zones" "all" {}

#
# Route53 hosted zone
#
data "aws_route53_zone" "apex" {
  name       = "${var.zoneapex}"
  depends_on = ["aws_route53_zone.apex"]
}

#
# ELB target used for S3 access
#
data "aws_elb_service_account" "aslsp" {}

#
# user data template and vars
#
data "template_file" "user_data" {
  template = "${file("user_data.yaml")}"

  vars {
    zoneapex = "${var.zoneapex}"
  }
}

#
# subnet data used with network ELB
#
data "aws_subnet_ids" "aslsp" {
  vpc_id = "vpc-a27c82c6"
}

data "aws_subnet" "aslsp" {
  count = "${length(data.aws_subnet_ids.aslsp.ids)}"
  id    = "${data.aws_subnet_ids.aslsp.ids[count.index]}"
}
