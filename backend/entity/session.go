package entity

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"seeloggyplus/backend/logger"
)

type Session struct {
	ID         string
	ServerInfo *Server
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

func (s *Session) Close() error {
	log := logger.Get()

	if s.SFTPClient != nil {
		if err := s.SFTPClient.Close(); err != nil {
			log.Warn().Msg("Failed Close SFTP Client")
		}
	}

	if s.SSHClient != nil {
		return s.SSHClient.Close()
	}
	return nil
}
