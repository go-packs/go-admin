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

## Try The Demo

This repository includes a runnable example application with seeded data:

```bash
go run ./examples
```

The example creates `examples/demo.db` on first run. Open `http://localhost:8080/admin` and sign in with:

```text
Email: admin@example.com
Password: password123
```

## Quick Start

```go
package main

import (
    "log"
    "net/http"

    admin "github.com/go-packs/go-admin"
    "github.com/go-packs/go-admin/server"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Product struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Price float64
}

func main() {
    db, _ := gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
    db.AutoMigrate(&Product{}, &admin.AdminUser{}, &admin.Permission{}, &admin.Session{}, &admin.AuditLog{})

    adm := admin.NewRegistry(db)

    adm.Register(Product{}).
        SetGroup("Inventory").
        RegisterField("ID", "ID", true).
        RegisterField("Name", "Product Name", false).
        RegisterField("Price", "Price", false)

    log.Println("Admin panel starting on http://localhost:8080/admin")
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

For full feature documentation including associations, scopes, custom actions, batch actions, charts, and custom pages, see the [Usage Guide](USAGE.md).

After GitHub Pages is enabled, the documentation site will be available at `https://go-packs.github.io/go-admin/`. Go API documentation is available on pkg.go.dev after the repository is public.

## License

This project is licensed under the MIT License.
