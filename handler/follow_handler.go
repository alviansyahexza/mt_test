package handler

import (
	"strconv"

	entity "github.com/alviansyahexza/mt_test/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) FollowUser(c *fiber.Ctx) error {
	user_id, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	follow := new(entity.Follow)
	if err := c.BodyParser(follow); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	f, err := h.followService.FollowUser(user_id, follow.FollowedId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	follow.Id = f.Id
	follow.FollowerId = f.FollowerId
	follow.CreatedAt = f.CreatedAt
	return c.Status(fiber.StatusCreated).JSON(follow)
}

func (h *Handler) UnfollowUser(c *fiber.Ctx) error {
	user_id, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	followed_id_str := c.Params("followed_id")
	followed_id, err := strconv.Atoi(followed_id_str)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid followed_id"})
	}
	h.followService.UnfollowUser(user_id, followed_id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Unfollowed user with ID: " + followed_id_str})
}
