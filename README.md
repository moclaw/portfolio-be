# Moclaw's Portfolio Backend API

REST API backend for my personal portfolio website. Built with Go and provides content management capabilities for the portfolio frontend.

## About This Project

This backend API serves my personal portfolio website, providing:

- **Content Management**: API endpoints for portfolio sections (projects, experience, technologies, etc.)
- **File Upload**: Handle image and file uploads for portfolio content
- **Authentication**: JWT-based admin authentication for content management
- **Database**: SQLite database for storing portfolio data

## Technology Stack

- **Language**: Go 1.24
- **Framework**: Gin Web Framework
- **Database**: SQLite with GORM ORM
- **Authentication**: JWT tokens
- **File Storage**: AWS S3 (LocalStack for development)
- **API Documentation**: Swagger/OpenAPI

## Frontend

The frontend application for this portfolio is available at: [https://github.com/Moclaw/portfolio](https://github.com/Moclaw/portfolio)

## 🚀 Getting Started

```bash
git clone https://github.com/Moclaw/portfolio-be
cd portfolio-be
go mod download
go run cmd/server/main.go
```

The API will be available at: `http://localhost:5303`

## 📊 API Documentation

Access the interactive Swagger documentation at: `http://localhost:5303/swagger/index.html`

## 📧 Contact

This backend powers my personal portfolio. Feel free to explore the API or check out the frontend application.

---

Moclaw's Portfolio Backend - Built with Go