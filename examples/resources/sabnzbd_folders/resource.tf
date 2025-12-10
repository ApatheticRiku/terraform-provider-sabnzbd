resource "sabnzbd_folders" "config" {
  download_dir           = "/data/incomplete"
  download_free          = "10G"
  complete_dir           = "/data/complete"
  complete_free          = "20G"
  auto_resume            = true
  permissions            = "755"
  watched_dir            = "/data/nzb-watch"
  watched_dir_scan_speed = 5
  scripts_dir            = "/config/scripts"
  password_file          = "/config/passwords.txt"
  nzb_backup_dir         = "/data/nzb-backup"
  admin_dir              = "/config/admin"
  backup_dir             = "/config/backup"
  log_dir                = "/config/logs"
}
