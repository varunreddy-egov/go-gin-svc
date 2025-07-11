package models

import (
	"github.com/google/uuid"
)

// TemplateConfigDB is the database model that matches the table schema
type TemplateConfigDB struct {
	ID               uuid.UUID         `gorm:"column:id;type:uuid;primary_key"`
	TemplateID       string            `gorm:"column:templateid;not null"`
	Version          string            `gorm:"column:version;not null"`
	TenantID         string            `gorm:"column:tenantid;not null"`
	FieldMapping     map[string]string `gorm:"column:fieldmapping;type:jsonb"`
	APIMapping       []APIMapping      `gorm:"column:apimapping;type:jsonb"`
	CreatedBy        string            `gorm:"column:createdby"`
	LastModifiedBy   string            `gorm:"column:lastmodifiedby"`
	CreatedTime      int64             `gorm:"column:createdtime"`
	LastModifiedTime int64             `gorm:"column:lastmodifiedtime"`
}

// ToDTO converts TemplateConfigDB to TemplateConfig (DB to API)
func (tc *TemplateConfigDB) ToDTO() TemplateConfig {
	return TemplateConfig{
		ID:           tc.ID,
		TemplateID:   tc.TemplateID,
		TenantID:     tc.TenantID,
		Version:      tc.Version,
		FieldMapping: tc.FieldMapping,
		APIMapping:   tc.APIMapping,
		AuditDetails: AuditDetails{
			CreatedBy:        tc.CreatedBy,
			CreatedTime:      tc.CreatedTime,
			LastModifiedBy:   tc.LastModifiedBy,
			LastModifiedTime: tc.LastModifiedTime,
		},
	}
}

// FromDTO converts TemplateConfig to TemplateConfigDB (API to DB)
func FromDTO(dto *TemplateConfig) TemplateConfigDB {
	return TemplateConfigDB{
		ID:               dto.ID,
		TemplateID:       dto.TemplateID,
		TenantID:         dto.TenantID,
		Version:          dto.Version,
		FieldMapping:     dto.FieldMapping,
		APIMapping:       dto.APIMapping,
		CreatedBy:        dto.AuditDetails.CreatedBy,
		CreatedTime:      dto.AuditDetails.CreatedTime,
		LastModifiedBy:   dto.AuditDetails.LastModifiedBy,
		LastModifiedTime: dto.AuditDetails.LastModifiedTime,
	}
}
