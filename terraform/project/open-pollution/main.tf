resource "null_resource" "add-ipfs-image" {
  provisioner "local-exec" {
    command = <<-EOT
      docker pull ipfs/go-ipfs:${local.ipfs-image-tag}
      docker tag ipfs/go-ipfs:${local.ipfs-image-tag} eu.gcr.io/open-pollution/go-ipfs:${local.ipfs-image-tag}
      docker push  eu.gcr.io/open-pollution/go-ipfs:${local.ipfs-image-tag}
    EOT
  }
}

module "ipfs-node" {
  source = "git::https://github.com/areknoster/public-distributed-commit-log.git//terraform/modules/gcp/gce_ipfs_node?ref=v1.3.0"

  ipfs-docker-image  = local.ipfs-image-tag
  registry-bucket-id = google_container_registry.registry.id
  subnetwork         = google_compute_subnetwork.open-pollution-subnetwork.self_link
  zone               = local.zone
  machine_type       = "g1-small"

  depends_on = [null_resource.add-ipfs-image]
}

module "sentinel" {
  source = "git::https://github.com/areknoster/public-distributed-commit-log.git//terraform/modules/gcp/gce_sentinel?ref=v1.3.0"

  ipfs-node-self-link = module.ipfs-node.ipfs-node-instance.self_link
  registry_bucket_id  = google_container_registry.registry.id
  sentinel_image      = "eu.gcr.io/${local.project}/${local.sentinel-image-name}:latest"
}
