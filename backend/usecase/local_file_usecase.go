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

		modTime := info.ModTime()
		fullPath := filepath.Join(path, info.Name())

		fileDTO := dto.FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: &modTime,
			Mode:    info.Mode().String(),
			Path:    fullPath,
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

func (l localFileUCImpl) GetDrives() ([]dto.DriveInfo, error) {
	var drives []dto.DriveInfo

	if runtime.GOOS == "windows" {
		// Get available drives on Windows
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			drivePath := string(drive) + ":\\"
			if _, err := os.Stat(drivePath); err == nil {
				// Get drive info
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
		// Unix-like systems - just return root
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
		return "", nil // Empty string means show drives
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
				ModTime: nil, // Will be handled safely in frontend
				Mode:    "drwxrwxrwx",
				IsDrive: true, // New field to identify drives
				Path:    drivePath,
			}
			files = append(files, fileDTO)
		}
	}

	return files, nil
}

func (l localFileUCImpl) getDriveLabel(drivePath string) string {
	// This is a simplified version - you might want to use Windows API for real drive labels
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
	// Simplified - you might want to detect removable, network drives, etc.
	return "local"
}
