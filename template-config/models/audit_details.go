package models

type AuditDetails struct {
	CreatedBy        string `json:"createdBy"`
	CreatedTime      int64  `json:"createdTime,string"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	LastModifiedTime int64  `json:"lastModifiedTime,string"`
}
