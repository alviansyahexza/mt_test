package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alviansyahexza/mt_test/config"
	"github.com/alviansyahexza/mt_test/repo"
	"github.com/alviansyahexza/mt_test/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	amqp "github.com/streadway/amqp"
)

type Handler struct {
	userService   *service.UserService
	postService   *service.PostService
	followService *service.FollowService
	jwt           *config.JWT
	db            *sql.DB
	redis         *redis.Client
	channel       *amqp.Channel
}

const CACHE_POST_QUEUE = "CACHE_POST_QUEUE"

func NewHandler(db *sql.DB, jwt *config.JWT, redis *redis.Client, channel *amqp.Channel) *Handler {
	handler := &Handler{
		userService:   service.NewUserService(repo.NewUserRepo(db)),
		postService:   service.NewPostService(db),
		followService: service.NewFollowService(db),
		jwt:           jwt,
		db:            db,
		redis:         redis,
		channel:       channel,
	}
	go consume(channel, handler)
	return handler
}

func (h *Handler) SendQueueJsonMessage(queueName string, message map[string]interface{}) error {
	if h.channel == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "RabbitMQ channel is not initialized")
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to marshal message to JSON: "+err.Error())
	}
	err = h.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send message to RabbitMQ: "+err.Error())
	}
	return nil
}

func consume(channel *amqp.Channel, handler *Handler) {
	_, err := channel.QueueDeclare(
		CACHE_POST_QUEUE, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		fmt.Println("Failed to declare a queue:", err)
		return
	} else {
		fmt.Println("Success declare: ", CACHE_POST_QUEUE)
	}
	msgs, err := channel.Consume(
		CACHE_POST_QUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("Failed to register a consumer: " + err.Error())
	}

	for msg := range msgs {
		fmt.Println(CACHE_POST_QUEUE, "-Received a message:", string(msg.Body))
		var message map[string]interface{}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			continue
		}
		postKey, ok := message["post_key"].(string)
		if !ok {
			fmt.Println("post_key is not a string")
			continue
		}
		userId, ok := message["user_id"].(float64)
		if !ok {
			fmt.Println("user_id is not a float64")
			continue
		}
		size, ok := message["size"].(float64)
		if !ok {
			fmt.Println("size is not a float64")
			continue
		}
		sortBy, ok := message["sort_by"].(string)
		if !ok {
			fmt.Println("sort_by is not a string")
			continue
		}
		sortOrder, ok := message["sort_order"].(string)
		if !ok {
			fmt.Println("sort_order is not a string")
			continue
		}
		isFeed, ok := message["is_feed"].(bool)
		if !ok {
			fmt.Println("is_feed is not a boolean")
			continue
		}

		handler.CachePost(context.Background(), postKey, int(userId), int(size), sortBy, sortOrder, isFeed)
	}
}

func (h *Handler) GetUserIdFromToken(c *fiber.Ctx) (int, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
	}
	tokenString := strings.Split(authHeader, " ")[1]
	t, err := h.jwt.ValidateToken(tokenString)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
	}
	userId := int(claims["user_id"].(float64))
	if userId == 0 {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "user not found")
	}
	return userId, nil
}
