package entity

import "time"

type Server struct {
	ID        string
	Name      string
	Address   string
	Port      int
	User      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
