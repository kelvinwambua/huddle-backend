
env "local" {
  src = "file://database/migrations"
  url = "postgres://postgres:King_kelvin1@localhost:5432/huddle?sslmode=disable"
  dev = "docker://postgres/16/dev"

  migration {
    dir = "file://database/migrations"
  }
}

env "docker" {
  src = "file://database/migrations"
  url = "postgres://user:password@postgres:5432/twitter_spaces?sslmode=disable"
  dev = "docker://postgres/16/dev"

  migration {
    dir = "file://database/migrations"
  }
}
