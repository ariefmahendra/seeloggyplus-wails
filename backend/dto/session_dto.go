package dto

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SessionManagerDto struct {
	ID         string                   `json:"id"`
	ServerInfo *ServerSessionManagement `json:"server_info"`
	SSHClient  *ssh.Client              `json:"-"`
	SFTPClient *sftp.Client             `json:"-"`
}
