package handler

import (
	"context"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/usecase"
)

type ConnectionHandler struct {
	connectionUC usecase.ConnectionUC
}

func NewConnectionHandler(uc usecase.ConnectionUC) *ConnectionHandler {
	return &ConnectionHandler{connectionUC: uc}
}

func (h *ConnectionHandler) TestConnection(ctx context.Context, payload *dto.ServerCreateRequest) error {
	log := logger.Get()
	log.Info().Str("address", payload.Address).Msg("TestConnection handler called")

	err := h.connectionUC.TestConnection(ctx, payload)
	if err != nil {
		log.Error().Err(err).Str("address", payload.Address).Msg("Connection test failed in handler")
		return err
	}

	log.Info().Str("address", payload.Address).Msg("Connection test successful in handler")
	return nil
}
