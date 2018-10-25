#
# hosted zone records
#
resource "aws_route53_record" "apex" {
  zone_id = "${data.aws_route53_zone.apex.zone_id}"
  name    = "${var.zoneapex}"
  type    = "A"

  alias {
    name                   = "${aws_lb.aslsp.dns_name}"
    zone_id                = "${aws_lb.aslsp.zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "wildcard" {
  zone_id = "${data.aws_route53_zone.apex.zone_id}"
  name    = "*.${var.zoneapex}"
  type    = "CNAME"
  ttl     = "60"
  records = ["${aws_lb.aslsp.dns_name}"]
}

resource "aws_route53_record" "www" {
  zone_id = "${data.aws_route53_zone.apex.zone_id}"
  name    = "www.${var.zoneapex}"
  type    = "CNAME"
  ttl     = "60"
  records = ["${aws_lb.aslsp.dns_name}"]
}
