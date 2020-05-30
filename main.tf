provider kcli {
  url = "127.0.0.1:50051"
}

resource "kcli_vm" "voila" {
  name = "voila"
  image = "fedora-coreos-31.20200420.3.0-qemu.x86_64.qcow2"
  overrides = "{'memory': 2048, 'numcpus': 3}"
}
