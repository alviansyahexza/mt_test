package handler

import (
	"strconv"

	entity "github.com/alviansyahexza/mt_test/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) FindPosts(c *fiber.Ctx) error {
	user_id, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil || size < 1 {
		size = 10
	}
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")
	if sortBy != "created_at" || (sortOrder != "asc" && sortOrder != "desc") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid sort parameters"})
	}
	postList, err := h.postService.GetPosts(user_id, page, size, sortBy, sortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(postList)
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	user_id, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	post := new(entity.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	p, err := h.postService.CreatePost(post.Title, post.Content, user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *Handler) GetPostById(c *fiber.Ctx) error {
	_, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}
	post, err := h.postService.GetPostById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
	}
	return c.JSON(post)
}

func (h *Handler) DeletePost(c *fiber.Ctx) error {
	user_id, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}
	err2 := h.postService.DeletePost(id, user_id)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err2.Error()})
	}
	return c.SendString("post deleted")
}

func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	user_id, err := h.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}
	post := new(entity.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	post.Id = id
	post.UserId = user_id
	updatedPost, err := h.postService.UpdatePost(post, user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedPost)
}
