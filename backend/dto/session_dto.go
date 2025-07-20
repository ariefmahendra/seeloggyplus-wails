package dto

import (
	"golang.org/x/crypto/ssh"
)

type SessionManagerDto struct {
	ID         string                   `json:"id"`
	ServerInfo *ServerSessionManagement `json:"serverInfo"`
	SSHClient  *ssh.Client              `json:"-"`
}
