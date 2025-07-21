package usecase

import (
	"os"
	"path/filepath"
	"runtime"
	"seeloggyplus/backend/dto"
	"sort"
	"strings"
)

type LocalFileUC interface {
	List(path string) ([]dto.FileInfo, error)
	GetUserHomeDir() (string, error)
	GetDrives() ([]dto.DriveInfo, error)
	GetRootPath() (string, error)
}

type localFileUCImpl struct {
}

func NewLocalFileUC() LocalFileUC {
	return &localFileUCImpl{}
}

func (l localFileUCImpl) List(path string) ([]dto.FileInfo, error) {
	if path == "" {
		if runtime.GOOS == "windows" {
			return l.listWindowsDrives()
		} else {
			path = "/"
		}
	}

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

		fullPath := filepath.Join(path, info.Name())
		fileDTO := dto.FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().String(),
			Mode:    info.Mode().String(),
			Path:    fullPath,
		}
		files = append(files, fileDTO)
	}

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

func (l localFileUCImpl) GetDrives() ([]dto.DriveInfo, error) {
	var drives []dto.DriveInfo

	if runtime.GOOS == "windows" {
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			drivePath := string(drive) + ":\\"
			if _, err := os.Stat(drivePath); err == nil {
				driveInfo := dto.DriveInfo{
					Name:   string(drive) + ":",
					Path:   drivePath,
					Label:  l.getDriveLabel(drivePath),
					Type:   l.getDriveType(drivePath),
					IsRoot: true,
				}
				drives = append(drives, driveInfo)
			}
		}
	} else {
		drives = append(drives, dto.DriveInfo{
			Name:   "Root",
			Path:   "/",
			Label:  "File System Root",
			Type:   "local",
			IsRoot: true,
		})
	}

	return drives, nil
}

func (l localFileUCImpl) GetRootPath() (string, error) {
	if runtime.GOOS == "windows" {
		return "", nil
	}
	return "/", nil
}

func (l localFileUCImpl) listWindowsDrives() ([]dto.FileInfo, error) {
	var files []dto.FileInfo

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		if _, err := os.Stat(drivePath); err == nil {
			// Create drive as directory entry
			fileDTO := dto.FileInfo{
				Name:    string(drive) + ":",
				Size:    0,
				IsDir:   true,
				ModTime: "",
				Mode:    "drwxrwxrwx",
				IsDrive: true,
				Path:    drivePath,
			}
			files = append(files, fileDTO)
		}
	}

	return files, nil
}

func (l localFileUCImpl) getDriveLabel(drivePath string) string {
	switch strings.ToUpper(drivePath[:1]) {
	case "C":
		return "Local Disk"
	case "D":
		return "Data"
	default:
		return "Local Disk"
	}
}

func (l localFileUCImpl) getDriveType(drivePath string) string {
	return "local"
}
