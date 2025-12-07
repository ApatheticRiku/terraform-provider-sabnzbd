# Configure a news server
resource "sabnzbd_server" "primary" {
  name        = "news.example.com"
  host        = "news.example.com"
  port        = 563
  username    = "myuser"
  password    = var.news_server_password
  connections = 20
  ssl         = true
  ssl_verify  = 2
  enable      = true
  priority    = 0
}

# Configure a backup/fill server
resource "sabnzbd_server" "backup" {
  name        = "backup.example.com"
  host        = "backup.example.com"
  port        = 563
  username    = "myuser"
  password    = var.backup_server_password
  connections = 10
  ssl         = true
  ssl_verify  = 2
  enable      = true
  optional    = true
  priority    = 1
}

variable "news_server_password" {
  description = "Password for primary news server"
  type        = string
  sensitive   = true
}

variable "backup_server_password" {
  description = "Password for backup news server"
  type        = string
  sensitive   = true
}
