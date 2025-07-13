package handler

import (
	"context"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/usecase"
)

type ServerHandler struct {
	serverUC usecase.ServerUseCase
}

func NewServerHandler(uc usecase.ServerUseCase) *ServerHandler {
	return &ServerHandler{serverUC: uc}
}

func (h *ServerHandler) AddServer(ctx context.Context, payload *dto.ServerCreateRequest) (*dto.ServerResponse, error) {
	log := logger.Get()
	log.Info().Str("name", payload.Name).Str("address", payload.Address).Msg("AddServer handler called")

	newServer, err := h.serverUC.AddServer(ctx, payload)
	if err != nil {
		log.Error().Err(err).Str("name", payload.Name).Msg("Failed to add server")
		return nil, err
	}

	log.Info().Str("id", newServer.ID).Str("name", newServer.Name).Msg("Server added successfully")
	return newServer, nil
}

func (h *ServerHandler) ListServers(ctx context.Context) ([]*dto.ServerResponse, error) {
	log := logger.Get()
	log.Info().Msg("ListServers handler called")

	servers, err := h.serverUC.ListServers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list servers")
		return nil, err
	}

	log.Info().Int("count", len(servers)).Msg("Servers listed successfully")
	return servers, nil
}

func (h *ServerHandler) UpdateServer(ctx context.Context, payload *dto.ServerUpdateRequest) error {
	log := logger.Get()
	log.Info().Str("id", payload.ID).Msg("UpdateServer handler called")

	err := h.serverUC.UpdateServer(ctx, payload)
	if err != nil {
		log.Error().Err(err).Str("id", payload.ID).Msg("Failed to update server")
		return err
	}

	log.Info().Str("id", payload.ID).Msg("Server updated successfully")
	return nil
}

func (h *ServerHandler) DeleteServer(ctx context.Context, id string) error {
	log := logger.Get()
	log.Info().Str("id", id).Msg("DeleteServer handler called")

	err := h.serverUC.RemoveServer(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete server")
		return err
	}

	log.Info().Str("id", id).Msg("Server deleted successfully")
	return nil
}

func (h *ServerHandler) GetServerById(ctx context.Context, id string) (*dto.ServerResponse, error) {
	log := logger.Get()
	log.Info().Str("id", id).Msg("GetServerById handler called")

	server, err := h.serverUC.GetBydId(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to get server by ID")
		return nil, err
	}

	log.Info().Str("id", id).Str("name", server.Name).Msg("Server retrieved successfully")
	return server, nil
}

func (h *ServerHandler) Migrate(ctx context.Context) error {
	log := logger.Get()
	log.Info().Msg("Migrate handler called")

	err := h.serverUC.Migrate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Database migration failed")
		return err
	}

	log.Info().Msg("Database migration completed successfully")
	return nil
}
