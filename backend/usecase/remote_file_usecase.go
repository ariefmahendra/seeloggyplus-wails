package usecase

import (
	"fmt"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/shared/util"
	"sort"
)

type RemoteFileUC interface {
	GetListFile(sessionID string, path string) ([]dto.FileInfo, error)
}

type remoteFIleUCImpl struct {
	sessionManagerUC SessionManagerUC
}

func NewRemoteFileUC(sessionManagerUC SessionManagerUC) RemoteFileUC {
	return &remoteFIleUCImpl{sessionManagerUC: sessionManagerUC}
}

func (r *remoteFIleUCImpl) GetListFile(sessionID string, path string) ([]dto.FileInfo, error) {
	log := logger.Get()

	session := r.sessionManagerUC.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("session with ID '%s' not found or has expired", sessionID)
	}

	log.Info().Str("path", path).Str("sessionID", sessionID).Msg("Listing files from remote server")

	sftpFiles, err := session.SFTPClient.ReadDir(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to list remote directory")
		return nil, util.NormalizeSSHConnectionError(err)
	}

	var files []dto.FileInfo
	for _, f := range sftpFiles {
		fileInfo := dto.FileInfo{
			Name:    f.Name(),
			Size:    f.Size(),
			IsDir:   f.IsDir(),
			ModTime: f.ModTime(),
			Mode:    f.Mode().String(),
		}
		files = append(files, fileInfo)
	}

	// Sorting: Directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}
