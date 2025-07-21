package usecase

import (
	"fmt"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/logger"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RemoteFileUC interface {
	GetListFile(sessionID string, path string) ([]dto.FileInfo, error)
}

type remoteFileUCImpl struct {
	sessionManagerUC SessionManagerUC
}

func NewRemoteFileUC(sessionManagerUC SessionManagerUC) RemoteFileUC {
	return &remoteFileUCImpl{sessionManagerUC: sessionManagerUC}
}

func (r *remoteFileUCImpl) GetListFile(sessionID string, path string) ([]dto.FileInfo, error) {
	log := logger.Get()

	session := r.sessionManagerUC.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("session with ID '%s' not found or has expired", sessionID)
	}

	if session.SSHClient == nil {
		return nil, fmt.Errorf("SSH client is not available for session '%s'", sessionID)
	}

	if strings.TrimSpace(path) == "" {
		path = "/"
	}

	log.Info().Str("path", path).Str("sessionID", sessionID).Msg("Listing files using SSH command")

	sshSession, err := session.SSHClient.NewSession()
	if err != nil {
		log.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to create SSH session")
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if closeErr := sshSession.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Str("sessionID", sessionID).Msg("Failed to close SSH session")
		}
	}()

	// Use ls with a specific format for better parsing
	cmd := fmt.Sprintf("ls -la '%s' 2>/dev/null || echo 'ERROR_DIR_NOT_ACCESSIBLE'", path)

	output, err := sshSession.Output(cmd)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to execute ls command")
		return nil, fmt.Errorf("failed to list directory '%s': %w", path, err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "ERROR_DIR_NOT_ACCESSIBLE") {
		return nil, fmt.Errorf("directory '%s' is not accessible", path)
	}

	files, err := r.parseListOutput(outputStr, path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to parse ls output")
		return nil, fmt.Errorf("failed to parse directory listing: %w", err)
	}

	// Sort: directories first, then alphabetically
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	log.Info().Str("path", path).Int("fileCount", len(files)).Msg("Successfully listed files")
	return files, nil
}

func (r *remoteFileUCImpl) parseListOutput(output, basePath string) ([]dto.FileInfo, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var files []dto.FileInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total ") {
			continue
		}

		if strings.HasSuffix(line, " .") || strings.HasSuffix(line, " ..") {
			continue
		}

		fileInfo, err := r.parseFileLine(line, basePath)
		if err != nil {
			log := logger.Get()
			log.Warn().Err(err).Str("line", line).Msg("Failed to parse file line")
			continue
		}

		if fileInfo != nil {
			files = append(files, *fileInfo)
		}
	}

	return files, nil
}

func (r *remoteFileUCImpl) parseFileLine(line, basePath string) (*dto.FileInfo, error) {
	fields := strings.Fields(line)
	if len(fields) < 9 {
		return nil, fmt.Errorf("invalid ls output format")
	}

	permissions := fields[0]
	isDir := strings.HasPrefix(permissions, "d")

	size, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		size = 0
	}

	dateStr := strings.Join(fields[5:8], " ")
	modTime := r.parseModTime(dateStr)

	filename := strings.Join(fields[8:], " ")
	fullPath := r.buildPath(basePath, filename)

	return &dto.FileInfo{
		Name:    filename,
		Size:    size,
		IsDir:   isDir,
		ModTime: modTime.String(),
		Mode:    permissions,
		Path:    fullPath,
	}, nil
}

func (r *remoteFileUCImpl) parseModTime(dateStr string) time.Time {
	timeFormats := []string{
		"Jan 2 15:04",  // Current year with time
		"Jan 2  2006",  // Different year
		"Jan _2 15:04", // Single digit day with time
		"Jan _2  2006", // Single digit day with year
	}

	for _, format := range timeFormats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			if parsedTime.Year() == 0 {
				now := time.Now()
				return time.Date(now.Year(), parsedTime.Month(), parsedTime.Day(),
					parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
					parsedTime.Nanosecond(), parsedTime.Location())
			}
			return parsedTime
		}
	}
	return time.Now() // Fallback to the current time
}

func (r *remoteFileUCImpl) buildPath(basePath, filename string) string {
	if basePath == "/" {
		return "/" + filename
	}
	return basePath + "/" + filename
}
