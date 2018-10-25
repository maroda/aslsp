resource "aws_iam_role" "aslsp" {
  name = "aslsp-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      }
    }
  ]
}
EOF
}

resource "aws_iam_policy" "aslsp" {
  name        = "aslsp-access"
  description = "S3 Access"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "arn:aws:s3:::${var.zoneapex}/*",
        "arn:aws:s3:::${var.zoneapex}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "aslsp" {
  name       = "aslsp-attachment"
  roles      = ["${aws_iam_role.aslsp.name}"]
  policy_arn = "${aws_iam_policy.aslsp.arn}"
}

resource "aws_iam_instance_profile" "aslsp" {
  name = "aslsp-profile"
  role = "${aws_iam_role.aslsp.name}"
}
