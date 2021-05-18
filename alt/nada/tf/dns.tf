resource "aws_route53_record" "in" {
  zone_id = "${var.zoneid}"
  name    = "in.${var.zoneapex}"
  type    = "CNAME"
  ttl     = 60
  records = ["${var.simple_lb}"]
}
