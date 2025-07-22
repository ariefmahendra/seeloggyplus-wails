package dto

type SettingsRequestDto struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SettingsResponseDto struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
