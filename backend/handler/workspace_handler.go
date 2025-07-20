package handler

import (
	"context"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/usecase"
)

type WorkspaceHandler struct {
	workspaceUC usecase.WorkspaceUC
}

func NewWorkspaceHandler(workspaceUC usecase.WorkspaceUC) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceUC: workspaceUC}
}

func (h *WorkspaceHandler) CreateWorkspace(ctx context.Context, req *dto.CreateWorkspaceRequest) (*dto.Workspace, error) {
	return h.workspaceUC.CreateWorkspace(ctx, req)
}

func (h *WorkspaceHandler) GetWorkspace(ctx context.Context, id string) (*dto.WorkspaceWithFiles, error) {
	return h.workspaceUC.GetWorkspace(ctx, id)
}

func (h *WorkspaceHandler) GetAllWorkspaces(ctx context.Context) ([]*dto.Workspace, error) {
	return h.workspaceUC.GetAllWorkspaces(ctx)
}

func (h *WorkspaceHandler) UpdateWorkspace(ctx context.Context, req *dto.UpdateWorkspaceRequest) error {
	return h.workspaceUC.UpdateWorkspace(ctx, req)
}

func (h *WorkspaceHandler) DeleteWorkspace(ctx context.Context, id string) error {
	return h.workspaceUC.DeleteWorkspace(ctx, id)
}

func (h *WorkspaceHandler) AddFilesToWorkspace(ctx context.Context, req *dto.AddFileToWorkspaceRequest) error {
	return h.workspaceUC.AddFilesToWorkspace(ctx, req)
}

func (h *WorkspaceHandler) RemoveFileFromWorkspace(ctx context.Context, fileID string) error {
	return h.workspaceUC.RemoveFileFromWorkspace(ctx, fileID)
}

func (h *WorkspaceHandler) AnalyzeFile(ctx context.Context, req *dto.AnalyzeFileRequest) (*dto.LogAnalysisResult, error) {
	return h.workspaceUC.AnalyzeFile(ctx, req)
}

func (h *WorkspaceHandler) Migrate(ctx context.Context) error {
	return h.workspaceUC.Migrate(ctx)
}
