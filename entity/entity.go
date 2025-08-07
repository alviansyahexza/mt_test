package entity

import "time"

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	Id         int       `json:"id"`
	FollowerId int       `json:"follower_id"`
	FollowedId int       `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}
