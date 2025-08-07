package main

import (
	"os"

	"github.com/alviansyahexza/mt_test/config"
	routes "github.com/alviansyahexza/mt_test/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load(".env.dev")
	}
	db := config.GetDbConnection()
	defer db.Close()
	jwt := config.NewJWT(os.Getenv("JWT_SECRET"))
	app := fiber.New()
	redis := config.GetRedisClient()
	defer redis.Close()
	channel, err := config.GetRabbitClient()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	routes.SetupRoutes(app, db, jwt, redis, channel, os.Getenv("JWT_SECRET"))
	app.Listen(":3000")
}
