resource "yandex_iam_service_account" "vm-sa" {
  folder_id = local.folder_id
  name      = local.sa_name
}

resource "yandex_resourcemanager_folder_iam_binding" "editor" {
  # Assign "editor" role to service account.
  folder_id = local.folder_id
  role      = "editor"
  members   = [
    "serviceAccount:${yandex_iam_service_account.vm-sa.id}"
  ]
}
