package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/services"
)

// Registry holds all the handlers for the API
type Registry struct {
	connectorHandler   *ConnectorHandler
	transformerHandler *TransformerHandler
	connectionHandler  *ConnectionHandler
}

// NewRegistry creates a new handler registry
func NewRegistry(
	connectorService *services.ConnectorService,
	transformerService *services.TransformerService,
	connectionService *services.ConnectionService,
) *Registry {
	return &Registry{
		connectorHandler:   NewConnectorHandler(connectorService),
		transformerHandler: NewTransformerHandler(transformerService),
		connectionHandler:  NewConnectionHandler(connectionService),
	}
}

// RegisterRoutes registers all API routes with the provided router
func (r *Registry) RegisterRoutes(router chi.Router) {
	router.Route("/v1", func(router chi.Router) {
		// Connectors
		router.Route("/connectors", func(router chi.Router) {
			router.Get("/", r.connectorHandler.ListConnectors)
			router.Post("/test", r.connectorHandler.TestConnector)
			router.Post("/fetch", r.connectorHandler.FetchData)
			router.Route("/{connectorID}", func(router chi.Router) {
				router.Get("/schema", r.connectorHandler.GetConnectorSchema)
			})
		})

		// Transformers
		router.Route("/transformers", func(router chi.Router) {
			router.Get("/", r.transformerHandler.ListTransformers)
			router.Post("/", r.transformerHandler.CreateTransformer)
			router.Post("/generate", r.transformerHandler.GenerateTransformer)
			router.Post("/transform", r.transformerHandler.TransformData)
			router.Route("/{transformerID}", func(router chi.Router) {
				router.Get("/", r.transformerHandler.GetTransformer)
				router.Put("/", r.transformerHandler.UpdateTransformer)
				router.Delete("/", r.transformerHandler.DeleteTransformer)
			})
		})

		// Connections
		router.Route("/connections", func(router chi.Router) {
			router.Get("/", r.connectionHandler.ListConnections)
			router.Post("/", r.connectionHandler.CreateConnection)
			router.Route("/{connectionID}", func(router chi.Router) {
				router.Get("/", r.connectionHandler.GetConnection)
				router.Put("/", r.connectionHandler.UpdateConnection)
				router.Delete("/", r.connectionHandler.DeleteConnection)
				router.Post("/execute", r.connectionHandler.ExecuteConnection)
				router.Patch("/active", r.connectionHandler.SetActive)
				router.Get("/executions", r.connectionHandler.GetExecutions)
			})
		})
	})
}
