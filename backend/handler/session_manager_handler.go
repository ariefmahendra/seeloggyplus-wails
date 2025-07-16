package handler

import (
	"context"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/usecase"
)

type SessionManagerHandler struct {
	sessionManagerUC usecase.SessionManagerUC
}

func NewSessionManagerHandler(sessionManagerUC usecase.SessionManagerUC) *SessionManagerHandler {
	return &SessionManagerHandler{sessionManagerUC: sessionManagerUC}
}

func (h *SessionManagerHandler) ConnectSession(ctx context.Context, server *dto.ServerSessionManagement) (string, error) {
	return h.sessionManagerUC.Connect(ctx, server)
}

func (h *SessionManagerHandler) GetSession(sessionID string) (*entity.Session, bool) {
	return h.sessionManagerUC.GetSession(sessionID)
}

func (h *SessionManagerHandler) CloseSession(sessionID string) error {
	return h.sessionManagerUC.CloseSession(sessionID)
}

func (h *SessionManagerHandler) CloseAllSession() {
	h.sessionManagerUC.CloseAllSession()
}
