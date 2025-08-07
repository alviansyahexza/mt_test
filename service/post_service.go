package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/alviansyahexza/mt_test/entity"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(title, content string, userId int) (*entity.Post, error) {
	post := &entity.Post{
		Title:   title,
		Content: content,
		UserId:  userId,
	}
	query := "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.QueryRow(query, title, content, userId).Scan(&post.Id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPosts(user_id int, page int, size int, sortBy string, sortOrder string, isFeed bool) ([]entity.Post, error) {
	posts := []entity.Post{}
	offset := (page - 1) * size
	var query string
	if isFeed {
		query = "SELECT id, title, content, user_id, created_at FROM posts WHERE user_id IN (SELECT followed_id FROM follows WHERE follower_id = $1) OR user_id = $1 ORDER BY " + sortBy + " " + sortOrder + " LIMIT $2 OFFSET $3"
	} else {
		query = "SELECT id, title, content, user_id, created_at FROM posts WHERE user_id = $1 ORDER BY " + sortBy + " " + sortOrder + " LIMIT $2 OFFSET $3"
	}
	rows, err := s.db.Query(query, user_id, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostService) GetPostByIds(ids []string, sortBy string, sortOrder string) ([]entity.Post, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`
		SELECT id, title, content, user_id, created_at
		FROM posts
		WHERE id IN (%s)
		ORDER BY %s %s
	`, strings.Join(placeholders, ","), sortBy, sortOrder)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostService) GetPostIds(userId int, size int, sortBy string, sortOrder string, isFeed bool) ([]int, error) {
	var query string
	if isFeed {
		query = "SELECT id FROM posts WHERE user_id IN (SELECT followed_id FROM follows WHERE follower_id = $1) OR user_id = $1 ORDER BY " + sortBy + " " + sortOrder + " LIMIT $2"
	} else {
		query = "SELECT id FROM posts WHERE user_id = $1 ORDER BY " + sortBy + " " + sortOrder + " LIMIT $2"
	}
	rows, err := s.db.Query(query, userId, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postsIds []int
	for rows.Next() {
		var post int
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}
		postsIds = append(postsIds, post)
	}

	return postsIds, nil
}

func (s *PostService) GetPostById(id int) (*entity.Post, error) {
	post := &entity.Post{}
	query := "SELECT id, title, content, user_id, created_at FROM posts WHERE id = $1"
	err := s.db.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) UpdatePost(post *entity.Post, user_id int) (*entity.Post, error) {
	postInDb, err2 := s.GetPostById(post.Id)
	if err2 != nil {
		return nil, err2
	}
	if postInDb.UserId != user_id {
		return nil, errors.New("unauthorized")
	}
	query := "UPDATE posts SET title = $1, content = $2 WHERE id = $3 RETURNING id"
	err := s.db.QueryRow(query, post.Title, post.Content, post.Id).Scan(&post.Id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) DeletePost(id int, user_id int) error {
	postInDb, err2 := s.GetPostById(id)
	if err2 != nil {
		return err2
	}
	if postInDb.UserId != user_id {
		return errors.New("unauthorized")
	}
	query := "DELETE FROM posts WHERE id = $1 AND user_id = $2 "
	_, err := s.db.Exec(query, id, user_id)
	if err != nil {
		return err
	}
	return nil
}
