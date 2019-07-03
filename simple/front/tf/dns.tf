resource "aws_route53_record" "app" {
  zone_id = "${var.zoneid}"
  name    = "app.${var.zoneapex}"
  type    = "CNAME"
  ttl     = 60
  records = ["${var.simple_lb}"]
}
