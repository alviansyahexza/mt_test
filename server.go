package main

import (
	"github.com/alviansyahexza/mt_test/config"
	routes "github.com/alviansyahexza/mt_test/routes"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	db := config.GetDbConnection()
	defer db.Close()
	jwtKey := "dummy_secret_key"
	jwt := config.NewJWT(jwtKey)
	app := fiber.New()
	redis := config.GetRedisClient()
	defer redis.Close()
	routes.SetupFreeRoutes(app, db, jwt)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtKey),
	}))
	routes.SetupRoutes(app, db, jwt, redis)
	app.Listen(":3000")
}
