package dto

type ApiConnectionCreateRequest struct {
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Type        string      `json:"type"`
	Configs     []ApiConfig `json:"configs"`
}

type ApiConfig struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
}
