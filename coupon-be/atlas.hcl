variable "envfile" {
  type    = string
  default = ".env"
}

locals {
  envfile = {
    for line in split("\n", file(var.envfile)) :
    trimspace(split("=", line)[0]) => regex("=\\s*\"?([^\"]*)\"?", line)[0]
    if !startswith(trimspace(line), "#") && length(split("=", line)) > 1
  }
  mysql_dsn = local.envfile.MYSQL_URL
  mysql_user = regex("^([^:]+):", local.mysql_dsn)[0]
  mysql_pass = regex("^[^:]+:([^@]+)@", local.mysql_dsn)[0]
  mysql_host = regex("@tcp\\(([^:]+):", local.mysql_dsn)[0]
  mysql_port = regex("@tcp\\([^:]+:([0-9]+)\\)", local.mysql_dsn)[0]
  mysql_db = regex("\\)/([^?]+)", local.mysql_dsn)[0]
  
  # Construct standard MySQL URL
  mysql_url = "mysql://${local.mysql_user}:${local.mysql_pass}@${local.mysql_host}:${local.mysql_port}/${local.mysql_db}"
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/model",
    "--dialect", "mysql",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url

  // Define the URL of the database which is managed
  // in this environment.
  url = "${local.mysql_url}"


  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "docker://mysql/latest/dev"  // Update Docker image for PostgreSQL

  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}