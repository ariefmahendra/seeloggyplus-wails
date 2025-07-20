package dto

import "time"

type Workspace struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type WorkspaceFile struct {
	ID          string    `json:"id" db:"id"`
	WorkspaceID string    `json:"workspace_id" db:"workspace_id"`
	FileName    string    `json:"file_name" db:"file_name"`
	FilePath    string    `json:"file_path" db:"file_path"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	SessionID   string    `json:"session_id" db:"session_id"` // "local" or session ID
	IsAnalyzed  bool      `json:"is_analyzed" db:"is_analyzed"`
	AddedAt     time.Time `json:"added_at" db:"added_at"`
}

type CreateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type UpdateWorkspaceRequest struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type AddFileToWorkspaceRequest struct {
	WorkspaceID string     `json:"workspace_id" validate:"required"`
	Files       []FileInfo `json:"files" validate:"required,min=1"`
	SessionID   string     `json:"session_id"` // "local" or session ID
	BasePath    string     `json:"base_path"`
}

type LogEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Line      int       `json:"line"`
}

type LogAnalysisResult struct {
	FileID     string     `json:"file_id"`
	FileName   string     `json:"file_name"`
	TotalLines int        `json:"total_lines"`
	LogEntries []LogEntry `json:"log_entries"`
	ErrorCount int        `json:"error_count"`
	WarnCount  int        `json:"warn_count"`
	InfoCount  int        `json:"info_count"`
	DebugCount int        `json:"debug_count"`
	AnalyzedAt time.Time  `json:"analyzed_at"`
}

type AnalyzeFileRequest struct {
	WorkspaceID string `json:"workspace_id" validate:"required"`
	FileID      string `json:"file_id" validate:"required"`
}

type WorkspaceWithFiles struct {
	Workspace
	Files []WorkspaceFile `json:"files"`
}
