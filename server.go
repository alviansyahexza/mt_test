package main

import (
	"github.com/alviansyahexza/mt_test/config"
	routes "github.com/alviansyahexza/mt_test/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db := config.GetDbConnection()
	defer db.Close()
	jwtKey := "dummy_secret_key"
	jwt := config.NewJWT(jwtKey)
	app := fiber.New()
	redis := config.GetRedisClient()
	defer redis.Close()
	channel, err := config.GetRabbitClient()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	routes.SetupRoutes(app, db, jwt, redis, channel, jwtKey)
	app.Listen(":3000")
}
