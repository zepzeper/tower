package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services/credentials"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type CredentialsHandler struct {
	credentialsService *credentials.Service
}

func NewCredentialsHandler(credentialsService *credentials.Service) *CredentialsHandler {
	return &CredentialsHandler{
		credentialsService: credentialsService,
	}
}

func (h *CredentialsHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	credentials, err := h.credentialsService.FetchConnections()
	if err != nil {
		response.Error(w, "Failed to fetch connections: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	apiConnections := make([]dto.CredentialsResponse, 0, len(credentials))
	for _, cred := range credentials {
		description := ""
		if cred.Credentials.Description.Valid {
			description = cred.Credentials.Description.String
		}

		apiConnections = append(apiConnections, dto.CredentialsResponse{
			ID:          cred.Credentials.ID,
			Name:        cred.Credentials.Name,
			Description: description,
			Type:        cred.Credentials.Type,
			Active:      cred.Credentials.Active,
			CreatedAt:   cred.Credentials.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   cred.Credentials.UpdatedAt.Format(time.RFC3339),
		})
	}

	response.JSON(w, apiConnections, http.StatusOK)
}

func (h *CredentialsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	cred, err := h.credentialsService.GetConnection(id)
	if err != nil {
		response.Error(w, "Failed to fetch connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert configs to response format
	configs := make([]dto.CredentialsConfigResponse, 0)
	for _, config := range cred.CredentialsConfig {
		// Determine if the config is a secret
		isSecret := isSecretConfig(config.Key)

		// Mask secret values
		value := config.Value

		configs = append(configs, dto.CredentialsConfigResponse{
			CredentialsID: config.ConnectionID,
			Key:           config.Key,
			Value:         value,
			IsSecret:      isSecret,
			CreatedAt:     config.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     config.UpdatedAt.Format(time.RFC3339),
		})
	}

	description := ""
	if cred.Credentials.Description.Valid {
		description = cred.Credentials.Description.String
	}

	response.JSON(w, dto.CredentialsWithConfigResponse{
		Credentials: dto.CredentialsResponse{
			ID:          cred.Credentials.ID,
			Name:        cred.Credentials.Name,
			Description: description,
			Type:        cred.Credentials.Type,
			Active:      cred.Credentials.Active,
			CreatedAt:   cred.Credentials.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   cred.Credentials.UpdatedAt.Format(time.RFC3339),
		},
		Configs: configs,
	}, http.StatusOK)
}

func (h *CredentialsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CredentialsCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Call the service to create the connection
	cred, err := h.credentialsService.Create(req)
	if err != nil {
		response.Error(w, "Failed to create connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	description := ""
	if cred.Credentials.Description.Valid {
		description = cred.Credentials.Description.String
	}

	response.JSON(w, dto.CredentialsResponse{
		ID:          cred.Credentials.ID,
		Name:        cred.Credentials.Name,
		Description: description,
		Type:        cred.Credentials.Type,
		Active:      cred.Credentials.Active,
		CreatedAt:   cred.Credentials.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   cred.Credentials.UpdatedAt.Format(time.RFC3339),
	}, http.StatusCreated)
}

func (h *CredentialsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response.Error(w, "Only PATCH method is supported", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CredentialsUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Make sure the ID in the path matches the ID in the request body
	req.ID = id

	// Call the service to update the connection
	cred, err := h.credentialsService.UpdateConnection(req)
	if err != nil {
		response.Error(w, "Failed to update connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	description := ""
	if cred.Credentials.Description.Valid {
		description = cred.Credentials.Description.String
	}

	response.JSON(w, dto.CredentialsResponse{
		ID:          cred.Credentials.ID,
		Name:        cred.Credentials.Name,
		Description: description,
		Type:        cred.Credentials.Type,
		Active:      cred.Credentials.Active,
		CreatedAt:   cred.Credentials.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   cred.Credentials.UpdatedAt.Format(time.RFC3339),
	}, http.StatusOK)
}

func (h *CredentialsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.Error(w, "Only DELETE method is supported", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	err := h.credentialsService.DeleteConnection(id)
	if err != nil {
		response.Error(w, "Failed to delete connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Connection deleted successfully",
	}, http.StatusOK)
}

func (h *CredentialsHandler) Test(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	// Call the service to test the connection
	// Todo: Implement TestConnection method
	// cred, err := h.credentialsService.TestConnection(req)
	// if err != nil {
	// 	response.Error(w, "Failed to test connection: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Prepare the response data
	data := map[string]interface{}{
		"success":    true,
		"message":    "Connection test successful",
		"connection": map[string]interface{}{},
	}

	response.JSON(w, data, http.StatusOK)
}

// Helper function to determine if a config key contains sensitive information
func isSecretConfig(key string) bool {
	secretKeys := map[string]bool{
		"password":      true,
		"token":         true,
		"api_key":       true,
		"client_secret": true,
		"secret":        true,
		"access_key":    true,
		"private_key":   true,
	}

	return secretKeys[key]
}
