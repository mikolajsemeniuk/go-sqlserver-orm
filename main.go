package main

import (
	"log"
	"server/sailing"

	// _ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	var config struct {
		// put your env variable here
	}

	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	dsn := "sqlserver://sa:P%40ssw0rd@localhost:1433?database=yachts"
	storage, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/sign-up", sailing.SignUp(storage.Where, storage.Create))
	app.Get("/sign-in", sailing.SignIn(storage.Where))
	app.Get("/sign-out", sailing.SignOut())
	app.Get("/authorize", sailing.Authorize())

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
