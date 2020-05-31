provider kcli {
  url = "127.0.0.1:50051"
}


//resource "kcli_vm" "zz" {
//  name = "zz"
//  image = "fedora-coreos-31.20200420.3.0-qemu.x86_64.qcow2"
//  overrides = "{'memory': 2048, 'numcpus': 3, 'nets': ['default','default']}"
//}

//resource "kcli_pool" "images2" {
// name = "images2"
// path = "/var/lib/libvirt/images2"
//}

resource "kcli_network" "froutos" {
 name = "froutos"
 cidr = "192.168.136.0/24"
}
