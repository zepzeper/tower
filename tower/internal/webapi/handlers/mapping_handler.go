package handlers

import (
	"net/http"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services"
)

type MappingHandler struct {
    mappingService *services.MappingService
}

func NewMappingHandler(mappingService *services.MappingService) *MappingHandler {
	return &MappingHandler{
		mappingService: mappingService,
	}
}


// Generate handles GET /api/mappings/schema?source={source}&target={target}
func (h *MappingHandler) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	source := r.URL.Query().Get("source")
	target := r.URL.Query().Get("target")

	if source == "" || target == "" {
		response.Error(w, "Missing 'source' or 'target' query parameters", http.StatusBadRequest)
		return
	}

	data, err := h.mappingService.GenerateMapping(source, target)
	if err != nil {
		response.Error(w, "Failed to generate mapping: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, data, http.StatusOK)
}
