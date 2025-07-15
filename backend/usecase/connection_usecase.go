package usecase

import (
	"context"
	"errors"
	"fmt"
	"net"
	"seeloggyplus/backend/custom_error"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
)

type ConnectionUseCase interface {
	TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error
}

type connectionUseCase struct{}

func NewConnectionUseCase() ConnectionUseCase {
	return &connectionUseCase{}
}

func (uc *connectionUseCase) TestConnection(ctx context.Context, req *dto.ServerCreateRequest) error {
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

func (uc *connectionUseCase) testSSHConnection(ctx context.Context, address, user, password string) error {
	log := logger.Get()

	type handshakeResult struct {
		conn ssh.Conn
		err  error
	}
	resultChan := make(chan handshakeResult, 1)

	// Launch the blocking operation in a separate goroutine.
	go func() {
		sshConfig := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		dialer := net.Dialer{Timeout: 5 * time.Second}
		conn, err := dialer.DialContext(ctx, "tcp", address)
		if err != nil {
			resultChan <- handshakeResult{nil, err}
			return
		}

		sshConn, _, _, err := ssh.NewClientConn(conn, address, sshConfig)

		select {
		case <-ctx.Done():
			log.Info().Msg("Goroutine detected context cancellation, abandoning result.")
			if sshConn != nil {
				if closeErr := sshConn.Close(); closeErr != nil {
					log.Warn().Err(closeErr).Msg("Error while closing SSH connection in goroutine")
				}
			}
			return // Exit silently.
		default:
			resultChan <- handshakeResult{sshConn, err}
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("Context cancelled by user, test aborted.")
		return uc.normalizeSSHConnectionError(ctx.Err())

	case result := <-resultChan:
		if result.err != nil {
			log.Error().Err(result.err).Str("address", address).Msg("SSH operation failed")
			return uc.normalizeSSHConnectionError(result.err)
		}

		defer func() {
			if closeErr := result.conn.Close(); closeErr != nil {
				log.Warn().Err(closeErr).Msg("Error while closing successful SSH connection")
			}
		}()
		log.Info().Str("address", address).Msg("SSH connection  test successful")
		return nil
	}
}

// normalizeSSHConnectionError translates technical SSH errors into professional, user-friendly messages.
func (uc *connectionUseCase) normalizeSSHConnectionError(err error) error {
	var userMessage string
	var opError *net.OpError

	switch {
	case errors.Is(err, context.Canceled):
		userMessage = "Connection Test Canceled."
	case errors.Is(err, context.DeadlineExceeded):
		userMessage = "Connection Timed Out: Server is unreachable or a firewall is blocking the connection."
	case errors.As(err, &opError):
		var syscallErr syscall.Errno
		if errors.As(opError.Err, &syscallErr) && errors.Is(syscallErr, syscall.ECONNREFUSED) {
			userMessage = "Connection Refused: Please check the host address and port."
		} else {
			userMessage = "Network Error: Could not connect to the server."
		}
	default:
		originalError := err.Error()
		if strings.Contains(originalError, "no common algorithm") {
			userMessage = "Connection Failed: Server uses incompatible security algorithms."
		} else if strings.Contains(originalError, "unable to authenticate") || strings.Contains(originalError, "permission denied") {
			userMessage = "Authentication Failed: Please check your username and password."
		} else {
			userMessage = "Connection Failed: An unexpected error occurred."
		}
	}

	return &custom_error.UserFacingError{
		UserMessage:   userMessage,
		InternalError: err,
	}
}
