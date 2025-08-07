package handler

import (
	"database/sql"
	"strings"

	"github.com/alviansyahexza/mt_test/config"
	"github.com/alviansyahexza/mt_test/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	userService   *service.UserService
	postService   *service.PostService
	followService *service.FollowService
	jwt           *config.JWT
	db            *sql.DB
}

func NewHandler(db *sql.DB, jwt *config.JWT) *Handler {
	return &Handler{
		userService:   service.NewUserService(db),
		postService:   service.NewPostService(db),
		followService: service.NewFollowService(db),
		jwt:           jwt,
		db:            db,
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
