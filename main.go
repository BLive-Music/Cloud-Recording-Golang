package main

import (
	"fmt"
	"log"

	"github.com/BLive-Music/Cloud-Recording-Golang/api"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func healthCheck(c *fiber.Ctx) error {
	fmt.Println("Health Check")
	return c.SendString("OK")
}

func ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func main() {
	// Set global configuration
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicln(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.AutomaticEnv()

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", healthCheck)
	app.Get("/ping", ping)
	api.MountRoutes(app)

	app.Listen(":" + viper.GetString("PORT"))

}
