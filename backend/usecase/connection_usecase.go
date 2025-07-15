package usecase

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"seeloggyplus/backend/custom_error"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/shared/util"
)

type ConnectionUseCase interface {
	TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error
}

type connectionUsecaseImpl struct{}

func NewConnectionUseCase() ConnectionUseCase {
	return &connectionUsecaseImpl{}
}

func (uc *connectionUsecaseImpl) TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error {
	log := logger.Get()
	log.Info().Str("address", req.Address).Int("port", req.Port).Str("type", req.Type).Msg("Attempting to test connection")

	connectionType := entity.ConnectionType(req.Type)
	if err := connectionType.IsValid(); err != nil {
		log.Warn().Err(err).Str("type", req.Type).Msg("Invalid connection type provided")
		return &custom_error.UserFacingError{
			UserMessage:   "Invalid Connection Type Provided.",
			InternalError: err,
		}
	}

	address := fmt.Sprintf("%s:%d", req.Address, req.Port)

	client, err := util.CreateSSHClient(ctx, address, req.User, req.Password)
	if err != nil {
		return err
	}

	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Error while closing SSH connection")
		}
	}(client)

	log.Info().Str("address", address).Msg("SSH connection test successful")
	return nil
}
