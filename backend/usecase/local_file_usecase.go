package usecase

import (
	"os"
	"seeloggyplus/backend/dto"
)

type LocalFileUC interface {
	List(path string) ([]dto.FileInfo, error)
	GetUserHomeDir() (string, error)
}

type localFileUCImpl struct {
}

func NewLocalFileUC() LocalFileUC {
	return &localFileUCImpl{}
}

func (l localFileUCImpl) List(path string) ([]dto.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Initialize with empty slice instead of nil
	files := []dto.FileInfo{}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileDTO := dto.FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime(),
			Mode:    info.Mode().String(),
		}
		files = append(files, fileDTO)
	}

	return files, nil
}

func (l localFileUCImpl) GetUserHomeDir() (string, error) {
	return os.UserHomeDir()
}
