provider kcli {
  url = "127.0.0.1:50051"
}

resource "kcli_vm" "hendrix" {
  name = "hendrix"
  image = "centos8"
  overrides = "{'memory': 2048, 'numcpus': 3}"
}
