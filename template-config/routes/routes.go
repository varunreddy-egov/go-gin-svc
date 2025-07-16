package routes

import (
	"template-config/config"
	"template-config/handlers"
	"template-config/repository"
	"template-config/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Initialize dependencies
	repo := repository.NewTemplateConfigRepository(db)
	svc := service.NewTemplateConfigService(repo)
	handler := handlers.NewTemplateConfigHandler(svc)

	// API routes
	api := router.Group("/template-config/v1")
	{
		// Template config management routes
		templateConfig := api.Group("/config")
		{
			templateConfig.POST("/", handler.CreateTemplateConfig)
			templateConfig.PUT("/", handler.UpdateTemplateConfig)
			templateConfig.GET("/", handler.SearchTemplateConfigs)
			templateConfig.DELETE("/", handler.DeleteTemplateConfig)
		}

		// Template config render route
		api.POST("/render", handler.RenderTemplateConfig)
	}

	return router
}
