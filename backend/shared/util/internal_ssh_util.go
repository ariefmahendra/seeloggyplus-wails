package util

import (
	"context"
	"errors"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"seeloggyplus/backend/custom_error"
	"seeloggyplus/backend/logger"
	"strings"
	"syscall"
	"time"
)

func CreateSSHClient(ctx context.Context, address, user, password string) (*ssh.Client, error) {
	log := logger.Get()

	dialer := net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("TCP Dial failed")
		return nil, NormalizeSSHConnectionError(err)
	}

	type handshakeResult struct {
		clientConn ssh.Conn
		channels   <-chan ssh.NewChannel
		requests   <-chan *ssh.Request
		err        error
	}

	// create channel for receive the result with buffer is 1
	resultChan := make(chan handshakeResult, 1)

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	go func() {
		c, channels, reqs, err := ssh.NewClientConn(conn, address, sshConfig)
		resultChan <- handshakeResult{c, channels, reqs, err}
	}()

	select {
	case <-ctx.Done(): // context canceled example: timeout or cancel by user
		err := conn.Close()
		if err != nil {
			log.Warn().Str("address", address).Msg("Failure Close TCP Connection")
		}
		log.Warn().Str("address", address).Msg("SSH connection cancelled by context")
		return nil, NormalizeSSHConnectionError(ctx.Err())

	case res := <-resultChan: // Goroutine is done and send the result
		if res.err != nil {
			// if handshake failure, close tcp connection
			err := conn.Close()
			if err != nil {
				log.Warn().Err(res.err).Str("address", address).Msg("Failure close tcp connection")
			}
			log.Error().Err(res.err).Str("address", address).Msg("SSH Handshake failed")
			return nil, NormalizeSSHConnectionError(res.err)
		}

		log.Info().Str("address", address).Msg("SSH handshake successful")
		return ssh.NewClient(res.clientConn, res.channels, res.requests), nil
	}
}

func NormalizeSSHConnectionError(err error) error {
	var userMessage string
	var opError *net.OpError

	switch {
	case errors.Is(err, context.Canceled):
		userMessage = "Operation canceled by the user."
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, os.ErrDeadlineExceeded):
		userMessage = "Operation timed out: The server is unreachable or a firewall is blocking the connection."
	case errors.As(err, &opError):
		var syscallErr syscall.Errno
		if errors.As(opError.Err, &syscallErr) && errors.Is(syscallErr, syscall.ECONNREFUSED) {
			userMessage = "Connection refused: Please check the host address and port."
		} else {
			userMessage = "Network error: Could not connect to the server."
		}
	case strings.Contains(err.Error(), "no such file or directory"):
		userMessage = "File not found: The specified path does not exist on the remote server."
	case strings.Contains(err.Error(), "permission denied"):
		userMessage = "Permission denied: You do not have the required permissions to access this resource."
	case strings.Contains(err.Error(), "connection reset by peer"):
		userMessage = "Connection reset: The server unexpectedly closed the connection."
	case strings.Contains(err.Error(), "broken pipe"):
		userMessage = "Connection error: Broken pipe, the connection was interrupted."
	case strings.Contains(err.Error(), "no common algorithm"):
		userMessage = "Connection failed: The server uses incompatible security algorithms."
	case strings.Contains(err.Error(), "unable to authenticate") || strings.Contains(err.Error(), "permission denied"):
		userMessage = "Authentication failed: Please check your username and password."
	case strings.Contains(err.Error(), "EOF"):
		userMessage = "Connection error: Unexpected end of file, the connection was closed."
	default:
		userMessage = "An unexpected error occurred: " + err.Error()
	}

	return &custom_error.UserFacingError{
		UserMessage:   userMessage,
		InternalError: err,
	}
}
