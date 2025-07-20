*# Portfolio Backend API

A modern Go-based REST API backend for portfolio management with file upload capabilities using SQLite and S3-compatible storage.

## Features

- **Content Management**: Full CRUD operations for portfolio content
- **File Upload**: Upload files to S3-compatible storage (LocalStack for development)
- **Database**: SQLite database with GORM ORM
- **API Documentation**: RESTful API with Swagger/OpenAPI documentation
- **Docker Support**: Containerized application with Docker Compose
- **CORS Support**: Cross-origin resource sharing enabled
- **Structured Logging**: Request logging and error handling
- **Pagination**: Built-in pagination for list endpoints
- **Interactive API Docs**: Swagger UI for testing and exploring APIs

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin Web Framework
- **Database**: SQLite with GORM ORM
- **File Storage**: AWS S3 (LocalStack for development)
- **API Documentation**: Swagger/OpenAPI with gin-swagger
- **Containerization**: Docker & Docker Compose

## Project Structure

```
portfolio-be/
├── server/
│   └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/           # HTTP request handlers
│   │   ├── middleware/         # Custom middlewares
│   │   └── router.go          # Route definitions
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── database/
│   │   └── sqlite.go          # Database initialization
│   ├── models/                # Data models
│   ├── repository/            # Data access layer
│   └── services/              # Business logic layer
├── pkg/
│   └── utils/                 # Utility functions
├── docs/                      # Auto-generated Swagger docs
├── docker-compose.yml         # Docker compose configuration
├── Dockerfile                 # Docker image configuration
├── Makefile                   # Build and development commands
├── .air.toml                  # Live reload configuration
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## API Endpoints

### Health Check

- `GET /health` - Health check endpoint

### API Documentation

- `GET /swagger/index.html` - Interactive Swagger UI documentation
- `GET /swagger/doc.json` - Swagger JSON specification
- `GET /swagger/swagger.yaml` - Swagger YAML specification

### Content Management

- `POST /api/v1/contents` - Create new content
- `GET /api/v1/contents` - Get all contents (with pagination)
- `GET /api/v1/contents/{id}` - Get content by ID
- `PUT /api/v1/contents/{id}` - Update content
- `DELETE /api/v1/contents/{id}` - Delete content

### File Upload

- `POST /api/v1/uploads` - Upload file
- `GET /api/v1/uploads` - Get all uploads (with pagination)
- `GET /api/v1/uploads/{id}` - Get upload by ID
- `DELETE /api/v1/uploads/{id}` - Delete upload

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized setup)

### Local Development

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Moclaw/portfolio-be
   cd portfolio-be
   ```
2. **Seed data**: 
   ```bash
   go run cmd/seed/main.go
   ```

3. **Run the application**:
   ```bash
   make dev
   ```

4. **The API will be available at**: `http://localhost:5303`

5. **Access Swagger Documentation**:
   - Interactive API Docs: `http://localhost:5303/swagger/index.html`
   - JSON Specification: `http://localhost:5303/swagger/doc.json`

### Using Make Commands

This project includes a Makefile with helpful commands:

```bash
# Install swag CLI tool
make install-swag

# Generate swagger documentation
make swagger

# Run the application
make run

# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Install dependencies
make deps
```

### Docker Setup

1. **Build and run with Docker Compose**:
   ```bash
   docker-compose up --build
   ```

2. **Access the application**:
   - API: `http://localhost:5303`
   - LocalStack S3: `http://localhost:4566`

## Environment Variables

| Variable               | Description               | Default                 |
| ---------------------- | ------------------------- | ----------------------- |
| `PORT`                 | Server port               | `5303`                  |
| `DATABASE_URL`         | SQLite database file path | `portfolio.db`          |
| `S3_ENDPOINT`          | S3 endpoint URL           | `http://localhost:4566` |
| `S3_REGION`            | S3 region                 | `us-east-1`             |
| `S3_BUCKET`            | S3 bucket name            | `portfolio-bucket`      |
| `S3_ACCESS_KEY_ID`     | S3 access key             | `test`                  |
| `S3_SECRET_ACCESS_KEY` | S3 secret key             | `test`                  |

## API Usage Examples

### Create Content
```bash
curl -X POST http://localhost:5303/api/v1/contents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Portfolio Project",
    "description": "A web application built with React",
    "body": "Detailed description of the project...",
    "category": "web-development",
    "tags": "react,javascript,portfolio",
    "status": "published"
  }'
```

### Upload File
```bash
curl -X POST http://localhost:5303/api/v1/uploads \
  -F "file=@/path/to/your/file.jpg"
```

### Get All Contents
```bash
curl "http://localhost:5303/api/v1/contents?page=1&limit=10&category=web-development"
```

## Development

### API Development with Swagger

This project uses Swagger/OpenAPI for API documentation. To work with the API docs:

1. **Generate/Regenerate Swagger docs** (do this after changing API comments):
   ```bash
   make swagger
   # or manually:
   # swag init -g server/main.go -o ./docs
   ```

2. **Access Interactive Documentation**:
   - Visit `http://localhost:5303/swagger/index.html` after starting the server
   - Test APIs directly from the browser interface

3. **Adding Swagger Comments**:
   - Add comments above your handler functions using Swagger annotations
   - Example:
   ```go
   // CreateContent godoc
   // @Summary Create a new content
   // @Description Create a new content item
   // @Tags content
   // @Accept json
   // @Produce json
   // @Param content body models.ContentRequest true "Content data"
   // @Success 201 {object} utils.Response
   // @Router /api/v1/contents [post]
   func (h *ContentHandler) CreateContent(c *gin.Context) {
       // handler implementation
   }
   ```

### Running Tests

```bash
go test ./...
```

### Building the Application

```bash
go build -o portfolio-be server/main.go
```

### Database Migrations

The application automatically runs database migrations on startup. The following tables are created:

- `contents` - Stores portfolio content items
- `uploads` - Stores file upload metadata

## Integration with Frontend

This backend is designed to work seamlessly with your portfolio frontend. The API provides:

- **CORS support** for cross-origin requests from your frontend
- **Structured JSON responses** with consistent format
- **File upload endpoints** that return direct URLs to uploaded files
- **Content management endpoints** for displaying portfolio items

Make sure to update the CORS middleware in `internal/api/middleware/cors.go` to include your frontend URL.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u github.com/aws/aws-sdk-go
go get -u github.com/gin-gonic/gin
```

### Bước 3: Tạo Cấu Trúc Dự Án

Tạo cấu trúc thư mục như sau:

```
portfolio-be/
├── main.go
├── models/
│   └── content.go
├── routes/
│   └── routes.go
└── services/
    └── s3_service.go
```

### Bước 4: Tạo Mô Hình Dữ Liệu

Trong `models/content.go`, định nghĩa mô hình cho nội dung:

```go
package models

import (
    "gorm.io/gorm"
)

type Content struct {
    ID      uint   `json:"id" gorm:"primaryKey"`
    Title   string `json:"title"`
    Body    string `json:"body"`
    FileURL string `json:"file_url"`
}

func Migrate(db *gorm.DB) {
    db.AutoMigrate(&Content{})
}
```

### Bước 5: Tạo Dịch Vụ S3

Trong `services/s3_service.go`, tạo dịch vụ để upload file lên S3:

```go
package services

import (
    "bytes"
    "context"
    "fmt"
    "mime/multipart"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
    S3Client *s3.S3
    Bucket    string
}

func NewS3Service() *S3Service {
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
        Endpoint: aws.String("http://localhost:4566"), // LocalStack endpoint
        S3ForcePathStyle: aws.Bool(true),
    }))
    return &S3Service{
        S3Client: s3.New(sess),
        Bucket:    "your-bucket-name",
    }
}

func (s *S3Service) UploadFile(file multipart.File, fileName string) (string, error) {
    buf := new(bytes.Buffer)
    buf.ReadFrom(file)
    fileBytes := buf.Bytes()

    _, err := s.S3Client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(s.Bucket),
        Key:    aws.String(fileName),
        Body:   bytes.NewReader(fileBytes),
        ContentType: aws.String("application/octet-stream"),
    })
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("http://localhost:4566/%s/%s", s.Bucket, fileName), nil
}
```

### Bước 6: Tạo API CRUD

Trong `routes/routes.go`, tạo các route cho API:

```go
package routes

import (
    "net/http"
    "portfolio-be/models"
    "portfolio-be/services"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, s3Service *services.S3Service) {
    router.POST("/content", func(c *gin.Context) {
        var content models.Content
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
            return
        }

        fileURL, err := s3Service.UploadFile(file, file.Filename)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
            return
        }

        content.Title = c.PostForm("title")
        content.Body = c.PostForm("body")
        content.FileURL = fileURL

        db.Create(&content)
        c.JSON(http.StatusOK, content)
    })

    router.GET("/content/:id", func(c *gin.Context) {
        var content models.Content
        id := c.Param("id")
        if err := db.First(&content, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
            return
        }
        c.JSON(http.StatusOK, content)
    })

    router.PUT("/content/:id", func(c *gin.Context) {
        var content models.Content
        id := c.Param("id")
        if err := db.First(&content, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
            return
        }

        c.BindJSON(&content)
        db.Save(&content)
        c.JSON(http.StatusOK, content)
    })

    router.DELETE("/content/:id", func(c *gin.Context) {
        id := c.Param("id")
        db.Delete(&models.Content{}, id)
        c.Status(http.StatusNoContent)
    })
}
```

### Bước 7: Tạo `main.go`

Trong `main.go`, khởi tạo ứng dụng:

```go
package main

import (
    "portfolio-be/models"
    "portfolio-be/routes"
    "portfolio-be/services"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    db, err := gorm.Open(sqlite.Open("portfolio.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    models.Migrate(db)

    s3Service := services.NewS3Service()

    router := gin.Default()
    routes.SetupRoutes(router, db, s3Service)

    router.Run(":5303")
}
```

### Bước 8: Chạy LocalStack

Để chạy LocalStack, bạn có thể sử dụng Docker. Chạy lệnh sau:

```bash
docker run -d -p 4566:4566 -p 4510-4559:4510-4559 localstack/localstack
```

### Bước 9: Chạy Ứng Dụng

Chạy ứng dụng Golang:

```bash
go run main.go
```

### Bước 10: Kiểm Tra API

Bạn có thể sử dụng Postman hoặc curl để kiểm tra các API CRUD:

- **Tạo nội dung**:
```bash
curl -X POST http://localhost:5303/content -F "title=My Title" -F "body=My Body" -F "file=@path/to/your/file"
```

- **Lấy nội dung**:
```bash
curl http://localhost:5303/content/1
```

- **Cập nhật nội dung**:
```bash
curl -X PUT http://localhost:5303/content/1 -d '{"title": "Updated Title", "body": "Updated Body"}'
```

- **Xóa nội dung**:
```bash
curl -X DELETE http://localhost:5303/content/1
```

### Kết Luận

Bạn đã hoàn thành việc tạo một dự án Golang với SQLite, tích hợp LocalStack S3 và API CRUD cho nội dung. Bạn có thể mở rộng thêm các tính năng khác theo nhu cầu của mình.
*