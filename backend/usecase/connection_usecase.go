package usecase

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/shared/util"
)

type ConnectionUC interface {
	TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error
}

type connectionUCImpl struct{}

func NewConnectionUseCase() ConnectionUC {
	return &connectionUCImpl{}
}

func (uc *connectionUCImpl) TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error {
	log := logger.Get()
	log.Info().Str("address", req.Address).Int("port", req.Port).Msg("Attempting to test connection")

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
