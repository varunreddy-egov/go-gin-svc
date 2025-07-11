package models

// RenderRequest represents the request for rendering
type RenderRequest struct {
	TemplateID string         `json:"templateId" binding:"required"`
	TenantID   string         `json:"tenantId" binding:"required"`
	Version    string         `json:"version" binding:"required"`
	Payload    map[string]any `json:"payload" binding:"required"`
}
