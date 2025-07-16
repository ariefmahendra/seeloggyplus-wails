package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/repository"
)

type ServerUC interface {
	AddServer(ctx context.Context, payload *dto.ServerCreateRequest) (*dto.ServerResponse, error)
	UpdateServer(ctx context.Context, payload *dto.ServerUpdateRequest) error
	RemoveServer(ctx context.Context, id string) error
	ListServers(ctx context.Context) ([]*dto.ServerResponse, error)
	GetBydId(ctx context.Context, id string) (*dto.ServerResponse, error)
	Migrate(ctx context.Context) error
}

type serverUCImpl struct {
	repo repository.ServerRepository
}

func NewServerUseCase(repo repository.ServerRepository) ServerUC {
	return &serverUCImpl{repo: repo}
}

func (s *serverUCImpl) Migrate(ctx context.Context) error {
	log := logger.Get()
	log.Info().Msg("Attempting to migrate server table")

	if err := s.repo.MigrateTable(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to migrate server table")
		return err
	}

	log.Info().Msg("Server table migration successful")
	return nil
}

func (s *serverUCImpl) AddServer(ctx context.Context, payload *dto.ServerCreateRequest) (*dto.ServerResponse, error) {
	log := logger.Get()
	serverID := uuid.New().String()
	log.Info().Str("name", payload.Name).Str("generated_id", serverID).Msg("Processing AddServer use case")

	connectionType := entity.ConnectionType(payload.Type)
	if err := connectionType.IsValid(); err != nil {
		log.Warn().Err(err).Str("type", payload.Type).Msg("Invalid connection type provided")
		return nil, err
	}

	newServer := entity.Server{
		ID:       serverID,
		Name:     payload.Name,
		Address:  payload.Address,
		Port:     payload.Port,
		User:     payload.User,
		Password: payload.Password,
		Type:     connectionType,
	}

	createdServer, err := s.repo.Create(ctx, &newServer)
	if err != nil {
		log.Error().Err(err).Str("name", payload.Name).Msg("Repository failed to create server")
		return nil, err
	}

	log.Info().Str("id", createdServer.ID).Msg("Server entity created successfully in repository")

	serverResponse := dto.ServerResponse{
		ID:        createdServer.ID,
		Name:      createdServer.Name,
		Address:   createdServer.Address,
		Port:      createdServer.Port,
		User:      createdServer.User,
		Password:  createdServer.Password,
		CreatedAt: createdServer.CreatedAt.String(),
		UpdatedAt: createdServer.UpdatedAt.String(),
		Type:      string(createdServer.Type),
	}

	return &serverResponse, nil
}

func (s *serverUCImpl) UpdateServer(ctx context.Context, payload *dto.ServerUpdateRequest) error {
	log := logger.Get()
	log.Info().Str("id", payload.ID).Msg("Processing UpdateServer use case")

	existingServer, err := s.repo.FindByID(ctx, payload.ID)
	if err != nil {
		log.Error().Err(err).Str("id", payload.ID).Msg("Failed to find server for update")
		return err
	}
	if existingServer == nil {
		err := fmt.Errorf("server with id %s not found", payload.ID)
		log.Warn().Err(err).Str("id", payload.ID).Msg("Attempted to update non-existent server")
		return err
	}

	connectionType := entity.ConnectionType(payload.Type)
	if err := connectionType.IsValid(); err != nil {
		log.Warn().Err(err).Str("type", payload.Type).Msg("Invalid connection type provided for update")
		return err
	}

	updatedServer := entity.Server{
		ID:       payload.ID,
		Name:     payload.Name,
		Address:  payload.Address,
		Port:     payload.Port,
		User:     payload.User,
		Password: payload.Password,
		Type:     connectionType,
	}

	if err := s.repo.Update(ctx, &updatedServer); err != nil {
		log.Error().Err(err).Str("id", payload.ID).Msg("Repository failed to update server")
		return err
	}

	log.Info().Str("id", payload.ID).Msg("Server entity updated successfully in repository")
	return nil
}

func (s *serverUCImpl) RemoveServer(ctx context.Context, id string) error {
	log := logger.Get()
	log.Info().Str("id", id).Msg("Processing RemoveServer use case")

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Repository failed to delete server")
		return err
	}

	log.Info().Str("id", id).Msg("Server entity removed successfully from repository")
	return nil
}

func (s *serverUCImpl) ListServers(ctx context.Context) ([]*dto.ServerResponse, error) {
	log := logger.Get()
	log.Info().Msg("Processing ListServers use case")

	servers, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Repository failed to find all servers")
		return nil, err
	}

	log.Info().Int("count", len(servers)).Msg("Successfully retrieved server entities from repository")

	var listServers []*dto.ServerResponse
	for _, server := range servers {
		listServers = append(listServers, &dto.ServerResponse{
			ID:        server.ID,
			Name:      server.Name,
			Address:   server.Address,
			Port:      server.Port,
			User:      server.User,
			Password:  server.Password,
			CreatedAt: server.CreatedAt.String(),
			UpdatedAt: server.UpdatedAt.String(),
			Type:      string(server.Type),
		})
	}

	return listServers, nil
}

func (s *serverUCImpl) GetBydId(ctx context.Context, id string) (*dto.ServerResponse, error) {
	log := logger.Get()
	log.Info().Str("id", id).Msg("Processing GetBydId use case")

	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Repository failed to find server by ID")
		return nil, err
	}
	if server == nil {
		err := fmt.Errorf("server with id %s not found", id)
		log.Warn().Err(err).Str("id", id).Msg("Server not found by ID")
		return nil, err
	}

	log.Info().Str("id", server.ID).Str("name", server.Name).Msg("Successfully retrieved server entity from repository")

	serverResponse := dto.ServerResponse{
		ID:        server.ID,
		Name:      server.Name,
		Address:   server.Address,
		Port:      server.Port,
		User:      server.User,
		Password:  server.Password,
		CreatedAt: server.CreatedAt.String(),
		UpdatedAt: server.UpdatedAt.String(),
		Type:      string(server.Type),
	}

	return &serverResponse, nil
}
