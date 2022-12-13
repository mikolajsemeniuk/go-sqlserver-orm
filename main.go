package main

import (
	"log"
	"server/sailing"

	_ "server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// @title           API
// @version         1.0
// @description     some API
// @BasePath /
// @schemes http https
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

	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	router.Post("/account/register", sailing.Register(storage.Where, storage.Create))
	router.Post("/account/login", sailing.Login(storage.Where))
	router.Get("/account/logout", sailing.Logout())
	router.Get("/account/authorize", sailing.Authorize())

	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Use(func(c *fiber.Ctx) error { return c.Status(fiber.StatusNotFound).Redirect("/swagger/index.html") })

	if err := router.Listen(":5000"); err != nil {
		log.Fatal(err)
	}
}
