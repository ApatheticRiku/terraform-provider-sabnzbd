# Configure the SABnzbd provider
provider "sabnzbd" {
  # URL of your SABnzbd instance
  url = var.sabnzbd_url

  # API key from SABnzbd Config -> General -> API Key
  # Can also be set via SABNZBD_API_KEY environment variable
  api_key = var.sabnzbd_api_key
}

variable "sabnzbd_api_key" {
  description = "SABnzbd API key"
  type        = string
  sensitive   = true
}

variable "sabnzbd_url" {
  description = "URL of your SABnzbd instance"
  type        = string
  default     = "http://localhost:8080"
}
