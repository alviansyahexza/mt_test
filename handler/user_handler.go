package handler

import (
	entity "github.com/alviansyahexza/mt_test/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SignUp(c *fiber.Ctx) error {
	user := new(entity.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	u, err := h.userService.SignUp(user.Name, user.Email, string(user.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(u)
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	user := new(entity.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	u, err := h.userService.SignIn(user.Email, user.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	s, err2 := h.jwt.GenerateToken(u.Id)
	if err2 != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "failed generate credentials"})
	}
	return c.SendString(s)
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	user := new(entity.User)
	userId, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	user.Id = userId
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	u, err2 := h.userService.UpdateProfile(user)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err2.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userId, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	u, err2 := h.userService.GetProfile(userId)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err2.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(u)
}
