package models

// TemplateConfigSearch represents search parameters
type TemplateConfigSearch struct {
	UUIDs      []string `form:"uuids"`
	TemplateID string   `form:"templateId"`
	TenantID   string   `form:"tenantId" binding:"required"`
	Version    string   `form:"version"`
}
