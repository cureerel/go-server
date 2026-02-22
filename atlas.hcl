env "local" {
  url = getenv("DATABASE_URL")
  dev = getenv("DEV_DATABASE_URL")
  migration {
    dir    = "file://migrations"
    format = atlas
  }
}

env "production" {
  url = getenv("DATABASE_URL")
  migration {
    dir    = "file://migrations"
    format = atlas
  }
}