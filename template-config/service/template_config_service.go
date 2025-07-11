package service

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"template-config/models"
	"template-config/repository"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type TemplateConfigService struct {
	repo       *repository.TemplateConfigRepository
	httpClient *resty.Client
}

func NewTemplateConfigService(repo *repository.TemplateConfigRepository) *TemplateConfigService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)

	return &TemplateConfigService{
		repo:       repo,
		httpClient: client,
	}
}

func (s *TemplateConfigService) Create(config *models.TemplateConfigDB) error {
	existing, err := s.repo.GetByTemplateIDAndVersion(config.TemplateID, config.TenantID, config.Version)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return fmt.Errorf("template config already exists for templateId: %s, tenantId: %s, version: %s",
			config.TemplateID, config.TenantID, config.Version)
	}

	config.ID = uuid.New()
	now := time.Now().Unix()
	config.CreatedTime = now
	config.LastModifiedTime = now
	// Optionally set CreatedBy/LastModifiedBy if you have user context
	// config.CreatedBy = ...
	// config.LastModifiedBy = ...

	return s.repo.Create(config)
}

func (s *TemplateConfigService) Update(config *models.TemplateConfigDB) error {
	existing, err := s.repo.GetByTemplateIDAndVersion(config.TemplateID, config.TenantID, config.Version)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	config.LastModifiedTime = now
	config.CreatedTime = existing.CreatedTime
	config.CreatedBy = existing.CreatedBy
	// Optionally set LastModifiedBy if you have user context
	// config.LastModifiedBy = ...

	return s.repo.Update(config)
}

func (s *TemplateConfigService) Search(search *models.TemplateConfigSearch) ([]models.TemplateConfigDB, error) {
	return s.repo.Search(search)
}

func (s *TemplateConfigService) Delete(templateID, tenantID, version string) error {
	// Check if template exists
	_, err := s.repo.GetByTemplateIDAndVersion(templateID, tenantID, version)
	if err != nil {
		return err
	}

	return s.repo.Delete(templateID, tenantID, version)
}

func (s *TemplateConfigService) Render(request *models.RenderRequest) (*models.RenderResponse, []models.Error) {
	// Get template config
	config, err := s.repo.GetByTemplateIDAndVersion(request.TemplateID, request.TenantID, request.Version)
	if err != nil {
		return nil, []models.Error{{
			Code:        "NOT_FOUND",
			Message:     "Template config not found",
			Description: err.Error(),
		}}
	}

	// Initialize response
	response := &models.RenderResponse{
		TemplateID: request.TemplateID,
		TenantID:   request.TenantID,
		Version:    request.Version,
		Data:       make(map[string]any),
	}

	// Apply field mappings
	for field, jsonPath := range config.FieldMapping {
		value := gjson.Get(fmt.Sprintf("%v", request.Payload), jsonPath)
		if value.Exists() {
			response.Data[field] = value.Value()
		}
	}

	// Execute API mappings in parallel
	if len(config.APIMapping) > 0 {
		errors := s.executeAPIMappings(config.APIMapping, request.Payload, response)
		if len(errors) > 0 {
			return nil, errors
		}
	}

	return response, nil
}

func (s *TemplateConfigService) executeAPIMappings(apiMappings []models.APIMapping, payload map[string]any, response *models.RenderResponse) []models.Error {
	var wg sync.WaitGroup
	errorChan := make(chan models.Error, len(apiMappings))

	for _, apiMapping := range apiMappings {
		wg.Add(1)
		go func(mapping models.APIMapping) {
			defer wg.Done()

			url := s.buildURL(mapping.Endpoint, payload)
			resp, err := s.httpClient.R().
				SetHeader("Content-Type", "application/json").
				Get(url)

			if err != nil {
				errorChan <- models.Error{
					Code:        "API_CALL_FAILED",
					Message:     "Failed to call external API",
					Description: err.Error(),
					Params:      []string{url, mapping.Method},
				}
				return
			}

			if resp.StatusCode() != http.StatusOK {
				errorChan <- models.Error{
					Code:        "API_CALL_FAILED",
					Message:     "External API returned non-200 status",
					Description: fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), resp.String()),
					Params:      []string{url, mapping.Method, fmt.Sprintf("%d", resp.StatusCode())},
				}
				return
			}

			responseData := gjson.Parse(resp.String())
			for field, jsonPath := range mapping.ResponseMapping {
				value := responseData.Get(jsonPath)
				if value.Exists() {
					response.Data[field] = value.Value()
				}
			}
		}(apiMapping)
	}

	wg.Wait()
	close(errorChan)

	var errors []models.Error
	for err := range errorChan {
		errors = append(errors, err)
	}
	return errors
}

func (s *TemplateConfigService) buildURL(endpoint models.EndpointConfig, payload map[string]any) string {
	url := endpoint.Base + endpoint.Path

	// Replace path parameters
	for param, jsonPath := range endpoint.PathParams {
		value := gjson.Get(fmt.Sprintf("%v", payload), jsonPath)
		if value.Exists() {
			url = strings.ReplaceAll(url, "{{"+param+"}}", fmt.Sprintf("%v", value.Value()))
		}
	}

	// Add query parameters
	if len(endpoint.QueryParams) > 0 {
		var queryParams []string
		for param, jsonPath := range endpoint.QueryParams {
			value := gjson.Get(fmt.Sprintf("%v", payload), jsonPath)
			if value.Exists() {
				queryParams = append(queryParams, fmt.Sprintf("%s=%v", param, value.Value()))
			}
		}
		if len(queryParams) > 0 {
			url += "?" + strings.Join(queryParams, "&")
		}
	}

	return url
}
