package dto

type CredentialsCreateRequest struct {
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	Type        string              `json:"type"`
	Configs     []CredentialsConfig `json:"configs"`
}

type CredentialsConfig struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
}

type CredentialsUpdateRequest struct {
	ID          string              `json:"id"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Active      *bool               `json:"active,omitempty"`
	Configs     []CredentialsConfig `json:"configs,omitempty"`
}

// Response DTOs
type CredentialsResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CredentialsConfigResponse struct {
	CredentialsID string `json:"connection_id"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	IsSecret      bool   `json:"is_secret"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CredentialsWithConfigResponse struct {
	Credentials CredentialsResponse         `json:"connection"`
	Configs     []CredentialsConfigResponse `json:"configs"`
}
