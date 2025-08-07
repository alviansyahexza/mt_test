package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	isFeed := c.Query("is_feed", "true") == "true"
	if sortBy != "created_at" || (sortOrder != "asc" && sortOrder != "desc") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid sort parameters"})
	}

	postKey := "post_" + strconv.Itoa(user_id) + "_" + sortBy + "_" + sortOrder + "_" + strconv.FormatBool(isFeed)
	idFromCache, errCache := h.findPostInCache(c.Context(), postKey, page, size)
	var postList []entity.Post
	if errCache != nil || len(idFromCache) == 0 {
		postList, err = h.postService.GetPosts(user_id, page, size, sortBy, sortOrder, isFeed)
		fmt.Println("Get post from database")
		h.SendQueueJsonMessage(CACHE_POST_QUEUE, map[string]interface{}{
			"post_key":   postKey,
			"user_id":    user_id,
			"size":       size,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
			"is_feed":    isFeed,
		})
		fmt.Println("Queued cache post: ", CACHE_POST_QUEUE)
	} else {
		postList, err = h.postService.GetPostByIds(idFromCache, sortBy, sortOrder)
		fmt.Println("Get post from cache")
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(postList)
}

func (h *Handler) findPostInCache(
	c context.Context,
	postKey string,
	page int,
	size int,
) ([]string, error) {

	sc, errorRedis := h.redis.Get(c, postKey).Result()
	if errorRedis != nil {
		return nil, errorRedis
	}

	s := strings.Split(sc, ",")
	if len(s) == 0 {
		return []string{}, nil
	}

	start := (page - 1) * size
	end := page * size
	if start > len(s) || end > len(s) {
		return nil, errors.New("Out of cached")
	}

	paged := s[start:end]
	return paged, nil
}

func (h *Handler) CachePost(context context.Context, postKey string, userId int, size int, sortBy string, sortOrder string, isFeed bool) {
	i, err := h.postService.GetPostIds(userId, 10, sortBy, sortOrder, isFeed)
	if err != nil || len(i) == 0 {
		fmt.Println("Failed to cache")
		return
	}

	idsString := ""
	for _, id := range i {
		idsString += (strconv.Itoa(id) + ",")
	}
	idsString = idsString[:len(idsString)-1]
	fmt.Println(idsString)
	_, errRedis := h.redis.Set(context, postKey, idsString, time.Second*10).Result()
	if errRedis != nil {
		fmt.Println(errRedis.Error())
	}
	fmt.Println("Done caching posts")
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
