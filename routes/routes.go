package routes

import (
	"database/sql"

	"github.com/alviansyahexza/mt_test/config"
	"github.com/alviansyahexza/mt_test/handler"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

func SetupRoutes(app *fiber.App, db *sql.DB, jwt *config.JWT, redis *redis.Client, rabbit *amqp.Channel, jwtKey string) {
	handler := handler.NewHandler(db, jwt, redis, rabbit)
	app.Post("/users", handler.SignUp)
	app.Post("/users/auth", handler.SignIn)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtKey),
	}))
	app.Get("/posts", handler.FindPosts)
	app.Post("/posts", handler.CreatePost)
	app.Get("/posts/:id", handler.GetPostById)
	app.Put("/posts/:id", handler.UpdatePost)
	app.Delete("/posts/:id", handler.DeletePost)
	app.Post("/users", handler.SignUp)
	app.Post("/users/auth", handler.SignIn)
	app.Put("/users", handler.UpdateUser)
	app.Post("/follows", handler.FollowUser)
	app.Delete("/follows/:followed_id", handler.UnfollowUser)
}
