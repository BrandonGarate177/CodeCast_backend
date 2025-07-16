env "local" {
  url = "postgres://user:pass@localhost:5432/codecast?sslmode=disable"

  migration {
    dir = "file://modules/anon_sessions/migrations"
  }
}