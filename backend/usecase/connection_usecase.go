package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"seeloggyplus/backend/custom_error"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
)

// ConnectionUseCase defines the interface for the connection use case.
type ConnectionUseCase interface {
	TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error
}

type connectionUseCase struct{}

// NewConnectionUseCase creates a new instance of ConnectionUseCase.
func NewConnectionUseCase() ConnectionUseCase {
	return &connectionUseCase{}
}

// TestConnection attempts to establish a connection to a server based on its type.
func (uc *connectionUseCase) TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error {
	log := logger.Get()
	log.Info().Str("address", req.Address).Int("port", req.Port).Str("type", req.Type).Msg("Attempting to test connection")

	connectionType := entity.ConnectionType(req.Type)
	if err := connectionType.IsValid(); err != nil {
		log.Warn().Err(err).Str("type", req.Type).Msg("Invalid connection type provided for testing")
		return &custom_error.UserFacingError{
			UserMessage:   "Invalid Connection Type Provided.",
			InternalError: err,
		}
	}

	address := fmt.Sprintf("%s:%d", req.Address, req.Port)

	switch connectionType {
	case entity.SFTP, entity.SCP:
		return uc.testSSHConnection(ctx, address, req.User, req.Password)
	default:
		internalErr := fmt.Errorf("connection type '%s' is not supported for testing", connectionType)

		return &custom_error.UserFacingError{
			UserMessage:   fmt.Sprintf("Connection Type '%s' Is Not Supported.", connectionType),
			InternalError: internalErr,
		}
	}
}

// testSSHConnection performs an SSH connection test.
func (uc *connectionUseCase) testSSHConnection(ctx context.Context, address, user, password string) error {
	log := logger.Get()

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Info().Str("user", user).Str("address", address).Msg("Dialing SSH server with context")

	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Failed to dial server")
		// Jika dial gagal, periksa apakah itu karena pembatalan.
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return uc.normalizeSSHConnectionError(err)
		}
		return uc.normalizeSSHConnectionError(err)
	}

	c, _, _, err := ssh.NewClientConn(conn, address, sshConfig)
	if err != nil {
		// ==================================================================
		// LOGIKA UTAMA UNTUK MENGATASI RACE CONDITION
		// ==================================================================
		// Periksa status context SEKARANG, setelah error terjadi.
		// ctx.Err() akan non-nil jika context telah dibatalkan.
		if ctx.Err() != nil {
			// Jika context dibatalkan, prioritaskan error pembatalan.
			log.Warn().Err(err).Msg("SSH handshake failed, but context was already canceled. Prioritizing cancellation error.")
			// Kirim error dari context agar dinormalisasi dengan benar.
			return uc.normalizeSSHConnectionError(ctx.Err())
		}
		// ==================================================================

		// Jika context TIDAK dibatalkan, maka ini adalah error handshake yang sah.
		log.Error().Err(err).Str("address", address).Msg("SSH handshake failed")
		return uc.normalizeSSHConnectionError(err)
	}
	defer c.Close()

	log.Info().Str("address", address).Msg("SSH connection test successful")
	return nil
}

// normalizeSSHConnectionError translates technical SSH custom_error into professional, user-friendly messages.
func (uc *connectionUseCase) normalizeSSHConnectionError(err error) error {
	originalError := err.Error()
	var userMessage string

	switch {
	// Menambahkan "context canceled" ke pemeriksaan ini adalah kunci.
	case strings.Contains(originalError, "i/o timeout") || strings.Contains(originalError, "context deadline exceeded") || strings.Contains(originalError, "context canceled"):
		userMessage = "Connection Timed Out or Canceled."

	case strings.Contains(originalError, "no common algorithm"):
		userMessage = "Connection Failed: Server uses incompatible security algorithms."
	case strings.Contains(originalError, "unable to authenticate") || strings.Contains(originalError, "permission denied"):
		userMessage = "Authentication Failed: Please check your username and password."
	case strings.Contains(originalError, "connection refused"):
		userMessage = "Connection Refused: Please check the host address and port."
	default:
		userMessage = "Connection Failed: An unexpected error occurred."
	}

	return &custom_error.UserFacingError{
		UserMessage:   userMessage,
		InternalError: err,
	}
}
