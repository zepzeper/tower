package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/database"
)

// Registry holds all the handlers for the API
type Registry struct {
	channelHandler     *ChannelHandler
	workflowHandler    *WorkflowHandler
	transformerHandler *TransformerHandler
	executionHandler   *ExecutionHandler
}

// NewRegistry creates a new handler registry
func NewRegistry(dbManager *db.Manager) *Registry {
	return &Registry{
		channelHandler:     NewChannelHandler(dbManager),
		workflowHandler:    NewWorkflowHandler(dbManager),
		transformerHandler: NewTransformerHandler(dbManager),
		executionHandler:   NewExecutionHandler(dbManager),
	}
}

// RegisterRoutes registers all API routes with the provided router
func (r *Registry) RegisterRoutes(router chi.Router) {
	// V1 API
	router.Route("/v1", func(router chi.Router) {
		// Channels
		router.Route("/channels", func(router chi.Router) {
			router.Get("/", r.channelHandler.ListChannels)
			router.Post("/", r.channelHandler.CreateChannel)
			router.Route("/{channelID}", func(router chi.Router) {
				router.Get("/", r.channelHandler.GetChannel)
				router.Put("/", r.channelHandler.UpdateChannel)
				router.Delete("/", r.channelHandler.DeleteChannel)
			})
		})

		// Workflows
		router.Route("/workflows", func(router chi.Router) {
			router.Get("/", r.workflowHandler.ListWorkflows)
			router.Post("/", r.workflowHandler.CreateWorkflow)
			router.Route("/{workflowID}", func(router chi.Router) {
				router.Get("/", r.workflowHandler.GetWorkflow)
				router.Put("/", r.workflowHandler.UpdateWorkflow)
				router.Delete("/", r.workflowHandler.DeleteWorkflow)
				router.Post("/execute", r.workflowHandler.ExecuteWorkflow)
				router.Patch("/toggle", r.workflowHandler.ToggleWorkflow)
			})
		})

		// Transformers
		router.Route("/transformers", func(router chi.Router) {
			router.Get("/", r.transformerHandler.ListTransformers)
			router.Post("/", r.transformerHandler.CreateTransformer)
			router.Route("/{transformerID}", func(router chi.Router) {
				router.Get("/", r.transformerHandler.GetTransformer)
				router.Put("/", r.transformerHandler.UpdateTransformer)
				router.Delete("/", r.transformerHandler.DeleteTransformer)
				router.Post("/test", r.transformerHandler.TestTransformer)
			})
		})

		// Executions
		router.Route("/executions", func(router chi.Router) {
			router.Get("/", r.executionHandler.ListExecutions)
			router.Get("/stats", r.executionHandler.GetExecutionStats)
			router.Route("/{executionID}", func(router chi.Router) {
				router.Get("/", r.executionHandler.GetExecution)
				router.Post("/cancel", r.executionHandler.CancelExecution)
			})
		})
	})
}
