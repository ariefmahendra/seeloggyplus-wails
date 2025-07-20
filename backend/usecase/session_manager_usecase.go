package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/repository"
	"seeloggyplus/backend/shared/util"
	"sync"
)

type SessionManagerUC interface {
	Connect(ctx context.Context, serverId string) (string, error)
	GetSession(sessionID string) *dto.SessionManagerDto
	CloseSession(sessionID string) error
	CloseAllSession()
	GetListSession() []*dto.SessionManagerDto
}

type sessionManagerUCImpl struct {
	sessions         map[string]*entity.Session
	mu               sync.RWMutex
	serverRepository repository.ServerRepository
}

func NewSessionManagerUC(serverRepository repository.ServerRepository) SessionManagerUC {
	return &sessionManagerUCImpl{
		sessions:         make(map[string]*entity.Session),
		serverRepository: serverRepository,
	}
}

func (uc *sessionManagerUCImpl) GetListSession() []*dto.SessionManagerDto {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	log := logger.Get()
	log.Info().Msg("Retrieving list of active SSH sessions")

	var sessionList []*dto.SessionManagerDto
	for _, session := range uc.sessions {
		// Always ensure an SSH client is available
		if session.SSHClient == nil {
			log.Warn().Str("sessionID", session.ID).Msg("Session has nil SSH client, skipping")
			continue
		}

		sessionManagerDto := &dto.SessionManagerDto{
			ID:         session.ID,
			ServerInfo: session.ServerInfo,
			SSHClient:  session.SSHClient,
		}
		sessionList = append(sessionList, sessionManagerDto)
	}

	log.Info().Int("count", len(sessionList)).Msg("Active SSH sessions retrieved")
	return sessionList
}

func (uc *sessionManagerUCImpl) Connect(ctx context.Context, serverId string) (string, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	log := logger.Get()
	log.Info().Str("serverId", serverId).Msg("Get Server by ID for SSH connection")

	server, err := uc.serverRepository.FindByID(ctx, serverId)
	if err != nil {
		log.Error().Err(err).Str("serverId", serverId).Msg("Failed to find server by ID")
		return "", fmt.Errorf("failed to find server with ID '%s': %w", serverId, err)
	}

	address := fmt.Sprintf("%s:%d", server.Address, server.Port)

	// Check for existing SSH connection by server ID
	for _, session := range uc.sessions {
		if session.ServerInfo.ID == server.ID {
			// Validate SSH client is still active
			if session.SSHClient != nil {
				// Test SSH connection health
				_, _, err := session.SSHClient.SendRequest("keepalive", false, nil)
				if err == nil {
					log.Info().Str("serverID", server.ID).Str("sessionID", session.ID).Msg("Reusing existing healthy SSH connection")
					return session.ID, nil
				}
				log.Warn().Str("serverID", server.ID).Str("sessionID", session.ID).Msg("Existing SSH connection is unhealthy, creating new one")
				// Clean up the unhealthy session
				delete(uc.sessions, session.ID)
				err = session.Close()
				if err != nil {
					log.Warn().Err(err).Str("serverID", server.ID).Str("sessionID", session.ID).Msg("Failed to close unhealthy SSH session")
				}
			}
		}
	}

	// Always create the SSH client first
	log.Info().Str("address", address).Str("user", server.User).Msg("Creating new SSH client connection")
	sshClient, err := util.CreateSSHClient(ctx, address, server.User, server.Password)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Failed to create SSH client")
		return "", err
	}

	// Verify the SSH client is working
	if sshClient == nil {
		log.Error().Str("address", address).Msg("SSH client is nil after creation")
		return "", fmt.Errorf("failed to establish SSH connection to %s", address)
	}

	log.Info().Str("address", address).Msg("SSH client created successfully")

	sessionID := uuid.New().String()

	serverInfo := dto.ServerSessionManagement{
		ID:       server.ID,
		Name:     server.Name,
		Address:  server.Address,
		Port:     server.Port,
		User:     server.User,
		Password: server.Password,
	}

	session := &entity.Session{
		ID:         sessionID,
		ServerInfo: &serverInfo,
		SSHClient:  sshClient,
	}

	uc.sessions[sessionID] = session

	log.Info().
		Str("sessionID", sessionID).
		Str("serverName", server.Name).
		Str("serverID", server.ID).
		Msg("New SSH session created successfully")

	return sessionID, nil
}

func (uc *sessionManagerUCImpl) GetSession(sessionID string) *dto.SessionManagerDto {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	log := logger.Get()
	session, found := uc.sessions[sessionID]
	if !found {
		log.Warn().Str("sessionID", sessionID).Msg("Session not found")
		return nil
	}

	// Always verify an SSH client exists
	if session.SSHClient == nil {
		log.Error().Str("sessionID", sessionID).Msg("Session exists but SSH client is nil")
		return nil
	}

	log.Info().Str("sessionID", sessionID).Msg("Retrieving session with active SSH client")

	sessionManagerDto := &dto.SessionManagerDto{
		ID:         session.ID,
		ServerInfo: session.ServerInfo,
		SSHClient:  session.SSHClient,
		// SFTPClient removed - focus on SSH only
	}

	return sessionManagerDto
}

func (uc *sessionManagerUCImpl) CloseSession(sessionID string) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	log := logger.Get()
	session, found := uc.sessions[sessionID]
	if !found {
		log.Warn().Str("sessionID", sessionID).Msg("Attempted to close non-existent session")
		return fmt.Errorf("session with ID '%s' not found", sessionID)
	}

	// Remove from a session map first
	delete(uc.sessions, sessionID)

	log.Info().Str("sessionID", sessionID).Msg("Closing SSH session")

	// Always ensure proper SSH client cleanup
	if session.SSHClient != nil {
		err := session.Close()
		if err != nil {
			log.Warn().Err(err).Str("sessionID", sessionID).Msg("Error occurred while closing SSH session")
			return err
		}
		log.Info().Str("sessionID", sessionID).Msg("SSH session closed successfully")
	} else {
		log.Warn().Str("sessionID", sessionID).Msg("Session had nil SSH client during close")
	}

	return nil
}

func (uc *sessionManagerUCImpl) CloseAllSession() {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	log := logger.Get()
	log.Info().Int("count", len(uc.sessions)).Msg("Closing all SSH sessions")

	for sessionID, session := range uc.sessions {
		log.Info().Str("sessionID", sessionID).Msg("Closing SSH session on app shutdown")

		if session.SSHClient != nil {
			err := session.Close()
			if err != nil {
				log.Warn().Err(err).Str("sessionID", sessionID).Msg("Failed to close SSH session during shutdown")
			} else {
				log.Info().Str("sessionID", sessionID).Msg("SSH session closed successfully during shutdown")
			}
		} else {
			log.Warn().Str("sessionID", sessionID).Msg("Session had nil SSH client during shutdown")
		}
	}

	// Clear all sessions
	uc.sessions = make(map[string]*entity.Session)
	log.Info().Msg("All SSH sessions closed and cleaned up")
}
