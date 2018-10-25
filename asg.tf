#
#	Launch Config
#
resource "aws_launch_configuration" "aslsp" {
  name_prefix          = "${var.zoneapex}-"
  image_id             = "${var.ami}"
  instance_type        = "${var.instance_type}"
  iam_instance_profile = "${aws_iam_instance_profile.aslsp.name}"
  security_groups      = ["${aws_security_group.aslsp.id}"]
  user_data            = "${data.template_file.user_data.rendered}"

  lifecycle {
    create_before_destroy = true
  }
}

#
#	ASG
#
resource "aws_autoscaling_group" "aslsp" {
  launch_configuration = "${aws_launch_configuration.aslsp.id}"
  availability_zones   = ["${data.aws_availability_zones.all.names}"]
  target_group_arns    = ["${aws_lb_target_group.aslsp-80.arn}", "${aws_lb_target_group.aslsp-443.arn}"]
  min_size             = "${var.asg_min}"
  max_size             = "${var.asg_max}"

  lifecycle {
    create_before_destroy = true
  }

  tag {
    key                 = "Name"
    value               = "aslsp"
    propagate_at_launch = true
  }
}
