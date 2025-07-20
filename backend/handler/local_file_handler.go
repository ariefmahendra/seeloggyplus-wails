package handler

import (
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/usecase"
)

type LocalFileHandler struct {
	localFileUC usecase.LocalFileUC
}

func NewLocalFileHandler(localFileUC usecase.LocalFileUC) *LocalFileHandler {
	return &LocalFileHandler{localFileUC: localFileUC}
}

func (h *LocalFileHandler) ListFiles(path string) ([]dto.FileInfo, error) {
	return h.localFileUC.List(path)
}

func (h *LocalFileHandler) GetUserHomeDir() (string, error) {
	return h.localFileUC.GetUserHomeDir()
}

func (h *LocalFileHandler) GetDrives() ([]dto.DriveInfo, error) {
	return h.localFileUC.GetDrives()
}

func (h *LocalFileHandler) GetRootPath() (string, error) {
	return h.localFileUC.GetRootPath()
}
