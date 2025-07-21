package dto

type FileInfo struct {
	SessionID string `json:"sessionID"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	IsDir     bool   `json:"isDir"`
	ModTime   string `json:"modTime"`
	Mode      string `json:"mode"`
	IsDrive   bool   `json:"isDrive"`
}

type DriveInfo struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Label  string `json:"label"`
	Type   string `json:"type"`
	IsRoot bool   `json:"isRoot"`
}
