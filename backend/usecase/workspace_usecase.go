package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/repository"
	"time"
)

type WorkspaceUC interface {
	CreateWorkspace(ctx context.Context, req *dto.CreateWorkspaceRequest) (*dto.Workspace, error)
	GetWorkspace(ctx context.Context, id string) (*dto.WorkspaceWithFiles, error)
	GetAllWorkspaces(ctx context.Context) ([]*dto.Workspace, error)
	UpdateWorkspace(ctx context.Context, req *dto.UpdateWorkspaceRequest) error
	DeleteWorkspace(ctx context.Context, id string) error

	AddFilesToWorkspace(ctx context.Context, req *dto.AddFileToWorkspaceRequest) error
	RemoveFileFromWorkspace(ctx context.Context, fileID string) error
	AnalyzeFile(ctx context.Context, req *dto.AnalyzeFileRequest) (*dto.LogAnalysisResult, error)

	Migrate(ctx context.Context) error
}

type workspaceUCImpl struct {
	workspaceRepo repository.WorkspaceRepository
	localFileUC   LocalFileUC
	remoteFileUC  RemoteFileUC
	sessionUC     SessionManagerUC
}

func NewWorkspaceUC(
	workspaceRepo repository.WorkspaceRepository,
	localFileUC LocalFileUC,
	remoteFileUC RemoteFileUC,
	sessionUC SessionManagerUC,
) WorkspaceUC {
	return &workspaceUCImpl{
		workspaceRepo: workspaceRepo,
		localFileUC:   localFileUC,
		remoteFileUC:  remoteFileUC,
		sessionUC:     sessionUC,
	}
}

func (uc *workspaceUCImpl) CreateWorkspace(ctx context.Context, req *dto.CreateWorkspaceRequest) (*dto.Workspace, error) {
	workspace := &dto.Workspace{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := uc.workspaceRepo.Create(ctx, workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	return workspace, nil
}

func (uc *workspaceUCImpl) GetWorkspace(ctx context.Context, id string) (*dto.WorkspaceWithFiles, error) {
	workspace, err := uc.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}
	if workspace == nil {
		return nil, fmt.Errorf("workspace not found")
	}

	files, err := uc.workspaceRepo.GetWorkspaceFiles(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace files: %w", err)
	}

	workspaceFiles := make([]dto.WorkspaceFile, len(files))
	for i, file := range files {
		workspaceFiles[i] = *file
	}

	return &dto.WorkspaceWithFiles{
		Workspace: *workspace,
		Files:     workspaceFiles,
	}, nil
}

func (uc *workspaceUCImpl) GetAllWorkspaces(ctx context.Context) ([]*dto.Workspace, error) {
	return uc.workspaceRepo.GetAll(ctx)
}

func (uc *workspaceUCImpl) UpdateWorkspace(ctx context.Context, req *dto.UpdateWorkspaceRequest) error {
	workspace, err := uc.workspaceRepo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	if workspace == nil {
		return fmt.Errorf("workspace not found")
	}

	workspace.Name = req.Name
	workspace.Description = req.Description
	workspace.UpdatedAt = time.Now()

	return uc.workspaceRepo.Update(ctx, workspace)
}

func (uc *workspaceUCImpl) DeleteWorkspace(ctx context.Context, id string) error {
	return uc.workspaceRepo.Delete(ctx, id)
}

func (uc *workspaceUCImpl) AddFilesToWorkspace(ctx context.Context, req *dto.AddFileToWorkspaceRequest) error {
	// Validate workspace exists
	workspace, err := uc.workspaceRepo.GetByID(ctx, req.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	if workspace == nil {
		return fmt.Errorf("workspace not found")
	}

	// Add each file to workspace
	for _, file := range req.Files {
		if file.IsDir {
			continue // Skip directories
		}

		workspaceFile := &dto.WorkspaceFile{
			ID:          uuid.New().String(),
			WorkspaceID: req.WorkspaceID,
			FileName:    file.Name,
			FilePath:    fmt.Sprintf("%s/%s", req.BasePath, file.Name),
			FileSize:    file.Size,
			SessionID:   req.SessionID,
			IsAnalyzed:  false,
			AddedAt:     time.Now(),
		}

		err := uc.workspaceRepo.AddFileToWorkspace(ctx, workspaceFile)
		if err != nil {
			return fmt.Errorf("failed to add file %s to workspace: %w", file.Name, err)
		}
	}

	return nil
}

func (uc *workspaceUCImpl) RemoveFileFromWorkspace(ctx context.Context, fileID string) error {
	return uc.workspaceRepo.RemoveFileFromWorkspace(ctx, fileID)
}

func (uc *workspaceUCImpl) AnalyzeFile(ctx context.Context, req *dto.AnalyzeFileRequest) (*dto.LogAnalysisResult, error) {
	// Implementation untuk analisis file log
	// Ini akan membaca file dan mengekstrak log entries
	return &dto.LogAnalysisResult{
		FileID:     req.FileID,
		AnalyzedAt: time.Now(),
		// TODO: Implement actual log parsing logic
	}, nil
}

func (uc *workspaceUCImpl) Migrate(ctx context.Context) error {
	return uc.workspaceRepo.Migrate(ctx)
}
