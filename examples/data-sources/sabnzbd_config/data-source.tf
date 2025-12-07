# Read SABnzbd configuration information
data "sabnzbd_config" "current" {}

output "sabnzbd_version" {
  description = "The version of SABnzbd"
  value       = data.sabnzbd_config.current.version
}

output "available_categories" {
  description = "List of configured categories"
  value       = data.sabnzbd_config.current.categories
}

output "available_scripts" {
  description = "List of available post-processing scripts"
  value       = data.sabnzbd_config.current.scripts
}
