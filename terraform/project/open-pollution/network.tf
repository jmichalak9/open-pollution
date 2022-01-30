resource "google_compute_network" "network" {
  name                    = "open-pollution-network"
  auto_create_subnetworks = false
}
resource "google_compute_subnetwork" "open-pollution-subnetwork" {
  name          = "subnetwork-${local.region}"
  ip_cidr_range = "10.0.0.0/16"
  region        = local.region
  network       = google_compute_network.network.id
}

resource "google_compute_firewall" "allow-ssh" {
  name    = "allow-ssh"
  network = google_compute_network.network.name
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  target_service_accounts = [ ] # TODO
  source_ranges = ["0.0.0.0/0"]
}
