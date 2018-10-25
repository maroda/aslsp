output "ns-apex" {
  value = ["${aws_route53_zone.apex.name_servers}"]
}
