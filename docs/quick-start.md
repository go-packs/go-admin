# Quick Start

This page creates a small admin app with SQLite and one `Product` resource.

## Install

```bash
go get github.com/go-packs/go-admin
go get gorm.io/driver/sqlite gorm.io/gorm
```

## Create `main.go`

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
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price float64
}

func main() {
	db, err := gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&Product{},
		&admin.AdminUser{},
		&admin.Permission{},
		&admin.Session{},
		&admin.AuditLog{},
	)

	adm := admin.NewRegistry(db)
	adm.Register(Product{}).
		SetGroup("Inventory").
		RegisterField("ID", "ID", true).
		RegisterField("Name", "Product Name", false).
		RegisterField("Price", "Price", false).
		SetFieldType("Price", "number")

	http.Handle("/admin/", server.NewRouter(adm))
	log.Println("Admin panel running at http://localhost:8080/admin")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Run

```bash
go run .
```

Open:

```text
http://localhost:8080/admin
```

## Next Steps

- Add an admin user with `admin.AdminUser` and `SetPassword`.
- Register more resources with `adm.Register`.
- Add filters with scopes, custom behavior with actions, and dashboard metrics with charts.
