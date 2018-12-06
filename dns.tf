#
# hosted zone records
#
resource "aws_route53_record" "wildcard" {
  zone_id = "${data.aws_route53_zone.apex.zone_id}"
  name    = "*.${var.zoneapex}"
  type    = "CNAME"
  ttl     = "60"
  records = ["${var.bptarget}"]
}
