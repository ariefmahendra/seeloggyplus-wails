package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/shared/util"
	"sync"
)

type SessionManagerUsecase interface {
	Connect(ctx context.Context, server *entity.Server) (string, error)
	GetSession(sessionID string) (*entity.Session, bool)
	CloseSession(sessionID string) error
	CloseAllSession()
}

type sessionManagerUsecaseImpl struct {
	sessions map[string]*entity.Session
	mu       sync.RWMutex
}

func NewSessionManagerUsecase() SessionManagerUsecase {
	return &sessionManagerUsecaseImpl{
		sessions: make(map[string]*entity.Session),
	}
}

func (uc *sessionManagerUsecaseImpl) Connect(ctx context.Context, server *entity.Server) (string, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	log := logger.Get()
	address := fmt.Sprintf("%s:%d", server.Address, server.Port)

	for _, session := range uc.sessions {
		if session.ServerInfo.ID == server.ID {
			log.Warn().Str("serverID", server.ID).Str("serverName", server.Name).Msg("Connection to this server already exists.")
			return session.ID, nil
		}
	}

	sshClient, err := util.CreateSSHClient(ctx, address, server.User, server.Password)
	if err != nil {
		return "", err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		err := sshClient.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to close ssh client")
		}
		log.Error().Err(err).Msg("Failed to create SFTP client")
		return "", util.NormalizeSSHConnectionError(err)
	}

	sessionID := uuid.New().String()
	session := &entity.Session{
		ID:         sessionID,
		ServerInfo: server,
		SSHClient:  sshClient,
		SFTPClient: sftpClient,
	}

	uc.sessions[sessionID] = session

	log.Info().Str("sessionID", sessionID).Str("serverName", server.Name).Msg("New persistent session created")
	return sessionID, nil
}

func (uc *sessionManagerUsecaseImpl) GetSession(sessionID string) (*entity.Session, bool) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	session, found := uc.sessions[sessionID]
	return session, found
}

func (uc *sessionManagerUsecaseImpl) CloseSession(sessionID string) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	log := logger.Get()
	session, found := uc.sessions[sessionID]
	if !found {
		return fmt.Errorf("session with ID '%s' not found", sessionID)
	}

	delete(uc.sessions, sessionID)
	log.Info().Str("sessionID", sessionID).Msg("Closing session")
	return session.Close()
}

func (uc *sessionManagerUsecaseImpl) CloseAllSession() {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	log := logger.Get()

	for _, session := range uc.sessions {
		log.Info().Str("sessionID", session.ID).Msg("Closing session on app shutdown")
		err := session.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed close session")
		}
	}

	uc.sessions = make(map[string]*entity.Session)
}
