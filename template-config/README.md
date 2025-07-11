# Template Config Service

A Go-based microservice for managing template configurations and rendering data with external API enrichment.

## Features

- **Template Config Management**: CRUD operations for template configurations
- **Data Transformation**: Field mapping from JSON payloads
- **API Enrichment**: Parallel external API calls with response mapping
- **Error Handling**: Detailed error reporting for failed API calls
- **Multi-tenant Support**: Tenant-based configuration isolation

## Technology Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **HTTP Client**: Resty
- **JSON Processing**: GJSON

## Project Structure

```
template-config/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go               # Configuration management
├── models/
│   └── template_config.go      # Data models
├── repository/
│   └── template_config_repo.go # Database operations
├── service/
│   └── template_config_service.go # Business logic
├── handlers/
│   └── template_config_handler.go # HTTP handlers
├── routes/
│   └── routes.go               # Route definitions
├── db/
│   └── postgres.go             # Database connection
├── go.mod                      # Go module file
└── README.md                   # This file
```

## Setup Instructions

### 1. Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Git

### 2. Environment Variables

Create a `.env` file with the following variables:

```bash
# Server Configuration
SERVER_ADDRESS=:8080
SERVER_PORT=8080

# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/template_config?sslmode=disable

# Logging Configuration
LOG_LEVEL=info

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=*
```

### 3. Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE template_config;
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The service will start on `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### 1. Create Template Config
- **POST** `/template-config`
- **Description**: Create a new template configuration
- **Request Body**: TemplateConfig object
- **Response**: Created TemplateConfig (201) or Error (400, 409, 500)

#### 2. Update Template Config
- **PUT** `/template-config`
- **Description**: Update an existing template configuration
- **Request Body**: TemplateConfig object
- **Response**: Updated TemplateConfig (200) or Error (400, 404, 500)

#### 3. Search Template Configs
- **GET** `/template-config`
- **Description**: Search template configurations
- **Query Parameters**:
  - `tenantId` (required): Tenant identifier
  - `templateId` (optional): Template identifier
  - `version` (optional): Version string
  - `uuids` (optional): Comma-separated list of UUIDs
- **Response**: Array of TemplateConfig (200) or Error (400, 404, 500)

#### 4. Delete Template Config
- **DELETE** `/template-config`
- **Description**: Delete a template configuration
- **Query Parameters**:
  - `templateId` (required): Template identifier
  - `tenantId` (required): Tenant identifier
  - `version` (required): Version string
- **Response**: Success (200) or Error (400, 404, 500)

#### 5. Render Template Config
- **POST** `/template-config/render`
- **Description**: Render template with data enrichment
- **Request Body**: RenderRequest object
- **Response**: RenderResponse (200) or Error (400, 404, 422, 500)

## Data Models

### TemplateConfig
```json
{
  "uuid": "string",
  "templateId": "string",
  "version": "string",
  "tenantId": "string",
  "fieldMapping": {
    "fieldName": "jsonPath"
  },
  "apiMapping": [
    {
      "method": "GET",
      "endpoint": {
        "base": "https://api.example.com",
        "path": "/users/{{userId}}",
        "pathParams": {
          "userId": "$.payload.user.id"
        },
        "queryParams": {
          "role": "$.payload.user.role"
        }
      },
      "responseMapping": {
        "userStatus": "$.response.status",
        "userExperience": "$.response.user.experience"
      }
    }
  ],
  "auditDetails": {
    "createdBy": "string",
    "createdTime": "2023-01-01T00:00:00Z",
    "lastModifiedBy": "string",
    "lastModifiedTime": "2023-01-01T00:00:00Z"
  }
}
```

### RenderRequest
```json
{
  "templateId": "string",
  "tenantId": "string",
  "version": "string",
  "payload": {
    "user": {
      "id": "123",
      "name": "John Doe",
      "role": "admin"
    }
  }
}
```

### RenderResponse
```json
{
  "templateId": "string",
  "tenantId": "string",
  "version": "string",
  "data": {
    "name": "John Doe",
    "userStatus": "active",
    "userExperience": "expert"
  },
  "errors": [
    {
      "endpoint": "https://api.example.com/users/123",
      "method": "GET",
      "error": "Connection timeout",
      "status": 500
    }
  ]
}
```

## Parallel API Processing

The render endpoint processes external API calls in parallel using Go goroutines. If any API calls fail:

1. The service continues processing other API calls
2. Failed calls are reported in the `errors` array
3. The response includes both successful data mappings and error details
4. HTTP status 422 is returned if any API calls fail

## Error Handling

- **400 Bad Request**: Invalid request format or missing required fields
- **404 Not Found**: Template config not found
- **409 Conflict**: Template config already exists
- **422 Unprocessable Entity**: API enrichment failed (with error details)
- **500 Internal Server Error**: Server-side errors

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/template-config cmd/main.go
```

### Docker Support
```bash
docker build -t template-config .
docker run -p 8080:8080 template-config
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License. 