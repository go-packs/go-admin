# Go Admin

A high-performance, modern, and extensible administration framework for Go, inspired by Active Admin. 

Go Admin uses reflection and GORM to automatically generate a professional-grade back-office for your database models with minimal configuration.

## Features

- 🔐 **Secure Authentication**: Session-based login with bcrypt password hashing.
- 📂 **Resource Grouping**: Organize your models into logical categories.
- 📊 **Visual Dashboard**: Customizable charts (powered by Chart.js) and stat widgets.
- 🔍 **Powerful Filtering**: Predefined scopes (tabs) and dynamic search filters.
- ⛓️ **Associations**: Automatic handling of `HasMany` and `BelongsTo` relationships.
- 📝 **Audit Logging**: Full history of every Create, Update, and Delete action.
- 📦 **Batch Actions**: Perform operations on multiple records at once.
- 📥 **CSV Export**: Export filtered data directly to CSV.
- 🎨 **Decorators**: Customize how fields are rendered (Currency, Badges, etc.).
- 🚀 **Portable**: Everything (HTML/CSS/JS) is bundled into your binary using `go:embed`.

## Installation

```bash
go get github.com/go-packs/go-admin
```

## Quick Start

```go
package main

import (
    "github.com/go-packs/go-admin"
    "github.com/go-packs/go-admin/server" // Import the server package for routing
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "net/http"
    "log"
)

type Product struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Price float64
}

func main() {
    db, _ := gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
    db.AutoMigrate(&Product{}, &admin.AdminUser{}, &admin.Permission{}, &admin.Session{}, &admin.AuditLog{})

    // Initialize Admin
    adm := admin.NewRegistry(db)

    // Register a Resource
    adm.Register(Product{}).
        SetGroup("Inventory").
        RegisterField("ID", "ID", true).
        RegisterField("Name", "Product Name", false).
        RegisterField("Price", "Price", false)

    // Start Server using the modular router
    log.Println("🚀 Admin panel starting on http://localhost:8080/admin")
    http.Handle("/admin/", server.NewRouter(adm))
    http.ListenAndServe(":8080", nil)
}
```

## Architecture & Project Structure

The project follows a modular architecture designed for maintainability and separation of concerns:

- `cmd/`: CLI tool for scaffolding and boilerplate generation.
- `config/`: Configuration management and defaults.
- `models/`: Core GORM models for users, sessions, and logs.
- `resource/`: Metadata definitions for administrative resources.
- `handlers/`: HTTP request handlers (Auth, CRUD, Export, etc.).
- `view/`: Template rendering and view logic.
- `server/`: Routing logic and HTTP middleware.
- `internal/`: Core business logic (Auth rules, Audit logging, CRUD services).
- `templates/`: HTML and CSS templates bundled via `go:embed`.

## Development

### Quality Control

We use `golangci-lint` for linting and `pre-commit` for local quality checks.

**Install pre-commit hooks:**
```bash
pre-commit install
```

**Run Linters:**
```bash
golangci-lint run
```

**Run Tests:**
```bash
go test ./...
```

## Documentation

For full feature documentation including **Associations**, **Scopes**, **Custom Actions**, and **Charts**, please refer to the [Usage Guide](USAGE.md).

## License

This project is licensed under the MIT License.
