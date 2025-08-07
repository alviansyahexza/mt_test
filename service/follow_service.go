package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alviansyahexza/mt_test/entity"
)

type FollowService struct {
	db *sql.DB
}

func NewFollowService(db *sql.DB) *FollowService {
	return &FollowService{db: db}
}

func (s *FollowService) FollowUser(followerId, followedId int) (*entity.Follow, error) {
	follow := &entity.Follow{
		FollowerId: followerId,
		FollowedId: followedId,
		CreatedAt:  time.Now(),
	}
	query := "INSERT INTO follows (follower_id, followed_id, created_at) VALUES ($1, $2, $3) RETURNING id"
	fmt.Println(followerId, followedId, follow.CreatedAt)
	err := s.db.QueryRow(query, followerId, followedId, follow.CreatedAt).Scan(&follow.Id)
	if err != nil {
		return nil, err
	}
	return follow, nil
}

func (s *FollowService) UnfollowUser(followerId, followedId int) error {
	query := "DELETE FROM follows WHERE follower_id = $1 AND followed_id = $2"
	_, err := s.db.Exec(query, followerId, followedId)
	return err
}

func (s *FollowService) GetFollowing(userId int) ([]int, error) {
	query := "SELECT followee_id FROM follows WHERE follower_id = $1"
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []int
	for rows.Next() {
		var followeeId int
		if err := rows.Scan(&followeeId); err != nil {
			return nil, err
		}
		following = append(following, followeeId)
	}
	return following, nil
}
