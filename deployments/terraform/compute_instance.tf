resource "yandex_compute_instance" "vm-1" {
  name               = "vm-1"
  zone               = "ru-central1-a"
  service_account_id = yandex_iam_service_account.vm-sa.id
  folder_id          = local.folder_id

  resources {
    cores         = 2
    memory        = 1
    core_fraction = 20
  }

  boot_disk {
    disk_id = yandex_compute_disk.vm-1-bd.id
  }

  network_interface {
    subnet_id      = yandex_vpc_subnet.subnet-a.id
    nat            = true
    nat_ip_address = yandex_vpc_address.vm-1-addr.external_ipv4_address[0].address
  }

  metadata = {
    user-data = file("meta.txt")
  }
}

resource "yandex_compute_disk" "vm-1-bd" {
  name      = "vm-1-bd"
  type      = "network-hdd"
  zone      = "ru-central1-a"
  size      = 15
  image_id  = data.yandex_compute_image.container-optimized-image.id
  folder_id = local.folder_id
}

data "yandex_compute_image" "container-optimized-image" {
  family = "container-optimized-image"
}

resource "yandex_vpc_address" "vm-1-addr" {
  name      = "vm-1-addr"
  folder_id = local.folder_id
  external_ipv4_address {
    zone_id = "ru-central1-a"
  }
}

output "external_ip_address_vm_1" {
  value = yandex_compute_instance.vm-1.network_interface.0.nat_ip_address
}
