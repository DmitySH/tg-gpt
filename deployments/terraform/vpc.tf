resource "yandex_vpc_network" "main-network" {
  folder_id   = local.folder_id
  description = "Network for the compute instance"
  name        = "main-network"
}

resource "yandex_vpc_subnet" "subnet-a" {
  folder_id      = local.folder_id
  description    = "Subnet in ru-central1-a availability zone"
  name           = "subnet-a"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.main-network.id
  v4_cidr_blocks = [local.zone_a_v4_cidr_blocks]
}
