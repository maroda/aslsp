##
##	Generate SSL certs
##

data "template_file" "certificates" {
  template   = "${file("${path.module}/template/kube-csr.json")}"
  depends_on = ["aws_elb.kube_api", "aws_instance.etcd", "aws_instance.controller", "aws_instance.worker"]

  vars {
    kube_api_elb_dns_name = "${aws_elb.kube_api.dns_name}"
    kube_cluster_dns      = "${var.kube_cluster_dns}"

    # Unfortunately, variables must be primitives, neither lists nor maps
    etcd0_ip       = "${aws_instance.etcd.0.private_ip}"
    etcd1_ip       = "${aws_instance.etcd.1.private_ip}"
    etcd2_ip       = "${aws_instance.etcd.2.private_ip}"
    controller0_ip = "${aws_instance.controller.0.private_ip}"
    controller1_ip = "${aws_instance.controller.1.private_ip}"
    controller2_ip = "${aws_instance.controller.2.private_ip}"
    worker0_ip     = "${aws_instance.worker.0.private_ip}"
    worker1_ip     = "${aws_instance.worker.1.private_ip}"
    worker2_ip     = "${aws_instance.worker.2.private_ip}"

    etcd0_dns       = "${aws_instance.etcd.0.private_dns}"
    etcd1_dns       = "${aws_instance.etcd.1.private_dns}"
    etcd2_dns       = "${aws_instance.etcd.2.private_dns}"
    controller0_dns = "${aws_instance.controller.0.private_dns}"
    controller1_dns = "${aws_instance.controller.1.private_dns}"
    controller2_dns = "${aws_instance.controller.2.private_dns}"
    worker0_dns     = "${aws_instance.worker.0.private_dns}"
    worker1_dns     = "${aws_instance.worker.1.private_dns}"
    worker2_dns     = "${aws_instance.worker.2.private_dns}"
  }
}

resource "null_resource" "certificates" {
  triggers {
    template_rendered = "${ data.template_file.certificates.rendered }"
  }

  provisioner "local-exec" {
    command = "echo '${ data.template_file.certificates.rendered }' > ../cert/kube-csr.json"
  }

  provisioner "local-exec" {
    command = "cd ../cert; cfssl gencert -initca ca-csr.json | cfssljson -bare ca; cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kube kube-csr.json | cfssljson -bare kube"
  }
}
