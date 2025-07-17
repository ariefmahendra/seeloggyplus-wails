package usecase

import (
	"os"
	"seeloggyplus/backend/dto"
	"sort"
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

	var files []dto.FileInfo
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

	// Sorting: Directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

func (l localFileUCImpl) GetUserHomeDir() (string, error) {
	return os.UserHomeDir()
}
