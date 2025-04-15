package models

// ChannelType defines the capabilities of a channel type
type ChannelType struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	IconURL     string            `json:"iconUrl"`
	AuthTypes   []AuthType        `json:"authTypes"`
	Triggers    []TriggerDef      `json:"triggers"`
	Actions     []ActionDef       `json:"actions"`
	ConfigDef   map[string]string `json:"configDef"`
}

// AuthType defines authentication methods
type AuthType struct {
	Type        string            `json:"type"` // oauth2, apiKey, basic, etc.
	DisplayName string            `json:"displayName"`
	Description string            `json:"description"`
	Fields      map[string]string `json:"fields"`
}

// TriggerDef defines a trigger that can be used with a channel
type TriggerDef struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	InputSchema map[string]SchemaField   `json:"inputSchema"`
	OutputSchema map[string]SchemaField  `json:"outputSchema"`
}

// ActionDef defines an action that can be performed with a channel
type ActionDef struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	InputSchema map[string]SchemaField   `json:"inputSchema"`
	OutputSchema map[string]SchemaField  `json:"outputSchema"`
}

// SchemaField defines a field in a data schema
type SchemaField struct {
	Type        string `json:"type"` // string, number, boolean, object, array
	Description string `json:"description"`
	Required    bool   `json:"required"`
}
