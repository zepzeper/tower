package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Response is a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Channel represents an API channel integration
type Channel struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Config      map[string]string `json:"config"`
}

// Workflow represents a data processing workflow
type Workflow struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Triggers    []Trigger    `json:"triggers"`
	Actions     []Action     `json:"actions"`
	Active      bool         `json:"active"`
}

// Trigger represents an event that can start a workflow
type Trigger struct {
	ChannelID string            `json:"channelId"`
	Event     string            `json:"event"`
	Config    map[string]string `json:"config"`
}

// Action represents an operation in a workflow
type Action struct {
	ChannelID     string            `json:"channelId"`
	Operation     string            `json:"operation"`
	Config        map[string]string `json:"config"`
	TransformerID string            `json:"transformerId,omitempty"`
}

// Transformer represents a data transformation configuration
type Transformer struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Mappings    []FieldMapping          `json:"mappings"`
	Functions   []TransformationFunction `json:"functions"`
}

// FieldMapping represents a mapping between fields
type FieldMapping struct {
	SourceField string `json:"sourceField"`
	TargetField string `json:"targetField"`
}

// TransformationFunction represents a function to apply to data
type TransformationFunction struct {
	Name       string   `json:"name"`
	TargetField string   `json:"targetField"`
	Parameters []string `json:"parameters"`
}

// Example data (in a real app, this would come from a database)
var channels = []Channel{
	{
		ID:          "ch_1",
		Name:        "Sample Twitter",
		Type:        "twitter",
		Description: "Twitter integration for social media",
		Config: map[string]string{
			"apiKey": "sample_key",
		},
	},
}

var workflows = []Workflow{
	{
		ID:          "wf_1",
		Name:        "Social to CRM",
		Description: "Send social media leads to CRM",
		Triggers: []Trigger{
			{
				ChannelID: "ch_1",
				Event:     "new_mention",
				Config:    map[string]string{},
			},
		},
		Actions: []Action{
			{
				ChannelID:     "ch_2",
				Operation:     "create_lead",
				TransformerID: "tr_1",
				Config:        map[string]string{},
			},
		},
		Active: true,
	},
}

var transformers = []Transformer{
	{
		ID:          "tr_1",
		Name:        "Twitter to CRM",
		Description: "Map Twitter data to CRM fields",
		Mappings: []FieldMapping{
			{
				SourceField: "user.screen_name",
				TargetField: "lead.twitter_handle",
			},
			{
				SourceField: "user.name",
				TargetField: "lead.name",
			},
		},
		Functions: []TransformationFunction{
			{
				Name:       "concatenate",
				TargetField: "lead.description",
				Parameters: []string{"'Lead from Twitter: '", "text"},
			},
		},
	},
}

// ----- Channel Handlers -----

// ListChannels returns all channels
func ListChannels(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    channels,
	})
}

// GetChannel returns a specific channel
func GetChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "channelID")
	for _, channel := range channels {
		if channel.ID == id {
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    channel,
			})
			return
		}
	}
	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Channel not found",
	})
}

// CreateChannel creates a new channel
func CreateChannel(w http.ResponseWriter, r *http.Request) {
	var channel Channel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// In a real app, would generate ID and save to database
	channel.ID = "ch_new"
	channels = append(channels, channel)

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    channel,
	})
}

// UpdateChannel updates an existing channel
func UpdateChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "channelID")
	
	var updated Channel
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	for i, channel := range channels {
		if channel.ID == id {
			updated.ID = id
			channels[i] = updated
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    updated,
			})
			return
		}
	}

	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Channel not found",
	})
}

// DeleteChannel deletes a channel
func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "channelID")
	
	for i, channel := range channels {
		if channel.ID == id {
			// Remove the channel (in a real app, might soft delete)
			channels = append(channels[:i], channels[i+1:]...)
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    nil,
			})
			return
		}
	}

	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Channel not found",
	})
}

// ----- Workflow Handlers -----

// ListWorkflows returns all workflows
func ListWorkflows(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    workflows,
	})
}

// GetWorkflow returns a specific workflow
func GetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "workflowID")
	for _, workflow := range workflows {
		if workflow.ID == id {
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    workflow,
			})
			return
		}
	}
	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Workflow not found",
	})
}

// CreateWorkflow creates a new workflow
func CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var workflow Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// In a real app, would generate ID and save to database
	workflow.ID = "wf_new"
	workflows = append(workflows, workflow)

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    workflow,
	})
}

// UpdateWorkflow updates an existing workflow
func UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "workflowID")
	
	var updated Workflow
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	for i, workflow := range workflows {
		if workflow.ID == id {
			updated.ID = id
			workflows[i] = updated
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    updated,
			})
			return
		}
	}

	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Workflow not found",
	})
}

// DeleteWorkflow deletes a workflow
func DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "workflowID")
	
	for i, workflow := range workflows {
		if workflow.ID == id {
			// Remove the workflow
			workflows = append(workflows[:i], workflows[i+1:]...)
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    nil,
			})
			return
		}
	}

	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Workflow not found",
	})
}

// ExecuteWorkflow manually executes a workflow
func ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "workflowID")
	
	var workflowFound bool
	for _, workflow := range workflows {
		if workflow.ID == id {
			workflowFound = true
			break
		}
	}

	if !workflowFound {
		respondWithJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Workflow not found",
		})
		return
	}

	// In a real app, would actually execute the workflow
	
	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"status": "workflow execution started"},
	})
}

// ----- Transformer Handlers -----

// ListTransformers returns all transformers
func ListTransformers(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    transformers,
	})
}

// GetTransformer returns a specific transformer
func GetTransformer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "transformerID")
	for _, transformer := range transformers {
		if transformer.ID == id {
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    transformer,
			})
			return
		}
	}
	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Transformer not found",
	})
}

// CreateTransformer creates a new transformer
func CreateTransformer(w http.ResponseWriter, r *http.Request) {
	var transformer Transformer
	if err := json.NewDecoder(r.Body).Decode(&transformer); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// In a real app, would generate ID and save to database
	transformer.ID = "tr_new"
	transformers = append(transformers, transformer)

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    transformer,
	})
}

// UpdateTransformer updates an existing transformer
func UpdateTransformer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "transformerID")
	
	var updated Transformer
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	for i, transformer := range transformers {
		if transformer.ID == id {
			updated.ID = id
			transformers[i] = updated
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    updated,
			})
			return
		}
	}

	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Transformer not found",
	})
}

// DeleteTransformer deletes a transformer
func DeleteTransformer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "transformerID")
	
	for i, transformer := range transformers {
		if transformer.ID == id {
			// Remove the transformer
			transformers = append(transformers[:i], transformers[i+1:]...)
			respondWithJSON(w, http.StatusOK, Response{
				Success: true,
				Data:    nil,
			})
			return
		}
	}

	respondWithJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   "Transformer not found",
	})
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"Failed to marshal JSON response"}`))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
