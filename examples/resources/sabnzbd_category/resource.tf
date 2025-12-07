# Configure download categories
resource "sabnzbd_category" "movies" {
  name     = "movies"
  dir      = "Movies"
  script   = "None"
  priority = 0
  pp       = "3" # +Repair/Unpack/Delete
}

resource "sabnzbd_category" "tv" {
  name     = "tv"
  dir      = "TV Shows"
  script   = "None"
  priority = 1
  pp       = "3"
}

resource "sabnzbd_category" "software" {
  name     = "software"
  dir      = "Software"
  script   = "None"
  priority = -1  # Low priority
  pp       = "2" # +Repair/Unpack
}
