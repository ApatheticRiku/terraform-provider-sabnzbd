# Terraform Provider for SABnzbd

This Terraform provider allows you to manage [SABnzbd](https://sabnzbd.org/) configuration as infrastructure. SABnzbd is a free and open-source Usenet binary newsreader.

## Features

- **News Servers** - Configure Usenet news servers with full SSL/TLS support
- **Categories** - Manage download categories with custom directories, scripts, and post-processing options
- **Folders** - Configure download paths, watched folders, scripts directory, and disk space management
- **Configuration Data** - Read SABnzbd version, available categories, and scripts

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for building)
- SABnzbd instance with API access enabled

## Installation

### From Source

```shell
git clone https://github.com/apatheticriku/terraform-provider-sabnzbd.git
cd terraform-provider-sabnzbd
go install
```

## Usage

### Provider Configuration

```hcl
provider "sabnzbd" {
  url     = "http://localhost:8080"
  api_key = var.sabnzbd_api_key
}
```

Or use environment variables:

```shell
export SABNZBD_URL="http://localhost:8080"
export SABNZBD_API_KEY="your-api-key"
```

### Example: Configure a News Server

```hcl
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
```

### Example: Configure Categories

```hcl
resource "sabnzbd_category" "movies" {
  name     = "movies"
  dir      = "Movies"
  priority = 0
  pp       = "3"  # +Repair/Unpack/Delete
}

resource "sabnzbd_category" "tv" {
  name     = "tv"
  dir      = "TV Shows"
  priority = 1
  pp       = "3"
}
```

### Example: Configure Folders

```hcl
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
}
```

### Example: Read Configuration

```hcl
data "sabnzbd_config" "current" {}

output "version" {
  value = data.sabnzbd_config.current.version
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `sabnzbd_server` | Manages news server configuration |
| `sabnzbd_category` | Manages download categories |
| `sabnzbd_folders` | Manages folder paths and disk space settings |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `sabnzbd_config` | Reads SABnzbd configuration (version, categories, scripts) |

## Development

### Building

```shell
make build
```

### Testing

```shell
make test        # Unit tests
make testacc     # Acceptance tests (requires running SABnzbd)
```

### Generating Documentation

```shell
make generate
```

## License

MPL-2.0
