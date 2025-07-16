package handler

import (
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/usecase"
)

type RemoteFileHandler struct {
	remoteFileUC usecase.RemoteFileUC
}

func NewRemoteFileHandler(remoteFileUC usecase.RemoteFileUC) *RemoteFileHandler {
	return &RemoteFileHandler{remoteFileUC: remoteFileUC}
}

func (h *RemoteFileHandler) GetListFiles(sessionID string, path string) ([]dto.FileInfo, error) {
	return h.remoteFileUC.GetListFile(sessionID, path)
}
