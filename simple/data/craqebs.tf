resource aws_ebs_volume "craqvol" {
  availability_zone = "us-east-2a"
  size              = 10
  type              = "gp2"

  tags = {
    Name = "craqvol"
  }
}
