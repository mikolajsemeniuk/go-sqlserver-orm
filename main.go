package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var config struct {
		// put your env variable here
	}

	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	storage, err := sql.Open("sqlserver", "server=localhost;user id=sa;password=P@ssw0rd;port=1433")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := storage.Exec("create database yachts"); err != nil {
		log.Fatal(err)
	}
}
