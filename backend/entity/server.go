package entity

import (
	"fmt"
	"time"
)

type Server struct {
	ID        string
	Name      string
	Address   string
	Port      int
	User      string
	Password  string
	Type      ConnectionType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConnectionType string

const (
	SFTP ConnectionType = "sftp"
	SCP  ConnectionType = "scp"
)

func (ct ConnectionType) IsValid() error {
	switch ct {
	case SFTP, SCP:
		return nil
	}
	return fmt.Errorf("invalid connection type: %s", string(ct))
}
