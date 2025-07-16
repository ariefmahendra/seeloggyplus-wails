package dto

import "time"

type FileInfo struct {
	SessionID string    `json:"sessionID"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"isDir"`
	ModTime   time.Time `json:"modTime"`
	Mode      string    `json:"mode"`
}
