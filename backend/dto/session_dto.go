package dto

type SessionManagerDto struct {
	ID         string                   `json:"id"`
	ServerInfo *ServerSessionManagement `json:"server_info"`
}
