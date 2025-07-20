package repository

import (
	"context"
	"database/sql"
	"fmt"
	"seeloggyplus/backend/dto"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *dto.Workspace) error
	GetByID(ctx context.Context, id string) (*dto.Workspace, error)
	GetAll(ctx context.Context) ([]*dto.Workspace, error)
	Update(ctx context.Context, workspace *dto.Workspace) error
	Delete(ctx context.Context, id string) error

	AddFileToWorkspace(ctx context.Context, file *dto.WorkspaceFile) error
	GetWorkspaceFiles(ctx context.Context, workspaceID string) ([]*dto.WorkspaceFile, error)
	RemoveFileFromWorkspace(ctx context.Context, fileID string) error
	UpdateFileAnalysisStatus(ctx context.Context, fileID string, isAnalyzed bool) error

	Migrate(ctx context.Context) error
}

type workspaceRepositoryImpl struct {
	db *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) WorkspaceRepository {
	return &workspaceRepositoryImpl{db: db}
}

func (r *workspaceRepositoryImpl) Create(ctx context.Context, workspace *dto.Workspace) error {
	query := `
        INSERT INTO workspaces (id, name, description, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
	_, err := r.db.ExecContext(ctx, query,
		workspace.ID, workspace.Name, workspace.Description,
		workspace.CreatedAt, workspace.UpdatedAt)
	return err
}

func (r *workspaceRepositoryImpl) GetByID(ctx context.Context, id string) (*dto.Workspace, error) {
	query := `
        SELECT id, name, description, created_at, updated_at
        FROM workspaces WHERE id = ?
    `
	workspace := &dto.Workspace{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&workspace.ID, &workspace.Name, &workspace.Description,
		&workspace.CreatedAt, &workspace.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return workspace, err
}

func (r *workspaceRepositoryImpl) GetAll(ctx context.Context) ([]*dto.Workspace, error) {
	query := `
        SELECT id, name, description, created_at, updated_at
        FROM workspaces ORDER BY created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*dto.Workspace
	for rows.Next() {
		workspace := &dto.Workspace{}
		err := rows.Scan(&workspace.ID, &workspace.Name, &workspace.Description,
			&workspace.CreatedAt, &workspace.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}
	return workspaces, rows.Err()
}

func (r *workspaceRepositoryImpl) Update(ctx context.Context, workspace *dto.Workspace) error {
	query := `
        UPDATE workspaces 
        SET name = ?, description = ?, updated_at = ?
        WHERE id = ?
    `
	_, err := r.db.ExecContext(ctx, query,
		workspace.Name, workspace.Description, workspace.UpdatedAt, workspace.ID)
	return err
}

func (r *workspaceRepositoryImpl) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete workspace files first
	_, err = tx.ExecContext(ctx, "DELETE FROM workspace_files WHERE workspace_id = ?", id)
	if err != nil {
		return err
	}

	// Delete workspace
	_, err = tx.ExecContext(ctx, "DELETE FROM workspaces WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *workspaceRepositoryImpl) AddFileToWorkspace(ctx context.Context, file *dto.WorkspaceFile) error {
	query := `
        INSERT INTO workspace_files (id, workspace_id, file_name, file_path, file_size, session_id, is_analyzed, added_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := r.db.ExecContext(ctx, query,
		file.ID, file.WorkspaceID, file.FileName, file.FilePath,
		file.FileSize, file.SessionID, file.IsAnalyzed, file.AddedAt)
	return err
}

func (r *workspaceRepositoryImpl) GetWorkspaceFiles(ctx context.Context, workspaceID string) ([]*dto.WorkspaceFile, error) {
	query := `
        SELECT id, workspace_id, file_name, file_path, file_size, session_id, is_analyzed, added_at
        FROM workspace_files WHERE workspace_id = ? ORDER BY added_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*dto.WorkspaceFile
	for rows.Next() {
		file := &dto.WorkspaceFile{}
		err := rows.Scan(&file.ID, &file.WorkspaceID, &file.FileName,
			&file.FilePath, &file.FileSize, &file.SessionID,
			&file.IsAnalyzed, &file.AddedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, rows.Err()
}

func (r *workspaceRepositoryImpl) RemoveFileFromWorkspace(ctx context.Context, fileID string) error {
	query := "DELETE FROM workspace_files WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, fileID)
	return err
}

func (r *workspaceRepositoryImpl) UpdateFileAnalysisStatus(ctx context.Context, fileID string, isAnalyzed bool) error {
	query := "UPDATE workspace_files SET is_analyzed = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, isAnalyzed, fileID)
	return err
}

func (r *workspaceRepositoryImpl) Migrate(ctx context.Context) error {
	// Create workspaces table
	createWorkspacesTable := `
    CREATE TABLE IF NOT EXISTS workspaces (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL
    )`

	// Create workspace_files table
	createWorkspaceFilesTable := `
    CREATE TABLE IF NOT EXISTS workspace_files (
        id TEXT PRIMARY KEY,
        workspace_id TEXT NOT NULL,
        file_name TEXT NOT NULL,
        file_path TEXT NOT NULL,
        file_size INTEGER NOT NULL,
        session_id TEXT,
        is_analyzed BOOLEAN DEFAULT FALSE,
        added_at DATETIME NOT NULL,
        FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
    )`

	if _, err := r.db.ExecContext(ctx, createWorkspacesTable); err != nil {
		return fmt.Errorf("failed to create workspaces table: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, createWorkspaceFilesTable); err != nil {
		return fmt.Errorf("failed to create workspace_files table: %w", err)
	}

	return nil
}
