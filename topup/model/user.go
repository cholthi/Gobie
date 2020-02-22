package model

import "time"

type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Token    string    `json:"token"`
	CreateAt time.Time `json:"created_at"`
}
