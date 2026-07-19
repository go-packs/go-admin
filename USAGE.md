# Usage Guide

This guide covers the main extension points in `github.com/go-packs/go-admin`. For a complete runnable app, see `examples/main.go`.

## Basic Setup

Create a GORM database connection, migrate your models and the built-in admin models, register resources, then mount the admin router.

```go
db, err := gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
if err != nil {
	log.Fatal(err)
}

db.AutoMigrate(&Product{}, &admin.AdminUser{}, &admin.Permission{}, &admin.Session{}, &admin.AuditLog{})

adm := admin.NewRegistry(db)
adm.Register(Product{}).
	SetGroup("Inventory").
	RegisterField("ID", "ID", true).
	RegisterField("Name", "Product Name", false).
	RegisterField("Price", "Price", false)

http.Handle("/admin/", server.NewRouter(adm))
log.Fatal(http.ListenAndServe(":8080", nil))
```

## Resources And Fields

Register a model with `Register`, then describe which fields should appear in the admin UI.

```go
adm.Register(Product{}).
	SetGroup("Products").
	RegisterField("ID", "ID", true).
	RegisterField("Name", "Product Name", false).
	RegisterField("Price", "Price", false).
	SetFieldType("Price", "number")
```

Use `SetIndexFields`, `SetShowFields`, and `SetEditFields` when each view needs a different field list.

## Field Types

Fields default to text inputs. You can customize a field with `SetFieldType`.

```go
adm.Register(User{}).
	RegisterField("Role", "Role", false).
	SetFieldType("Role", "select", "admin", "editor", "viewer")
```

Common field types include `text`, `number`, `select`, and `image`.

## Decorators

Decorators customize how values render in list and detail views.

```go
adm.Register(Product{}).
	RegisterField("Price", "Price", false).
	SetDecorator("Price", func(value interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("<strong>$%.2f</strong>", value.(float64)))
	})
```

## Associations

Use `HasMany` and `BelongsTo` to connect registered resources.

```go
adm.Register(Product{}).
	HasMany("ProductInfo", "Technical Specifications", "ProductInfo", "ProductID")

adm.Register(ProductInfo{}).
	BelongsTo("ProductID", "Parent Product", "Product", "ID")
```

## Scopes

Scopes add predefined filters to resource list pages.

```go
adm.Register(User{}).
	AddScope("admins", "Admins", func(db *gorm.DB) *gorm.DB {
		return db.Where("role = ?", "admin")
	})
```

## Actions

Member actions run against one record. Collection actions run from the resource page.

```go
adm.Register(Product{}).
	AddCollectionAction("discount", "Apply 10% Discount", func(res *admin.Resource, w http.ResponseWriter, r *http.Request) {
		db.Model(&Product{}).Where("price > ?", 0).Update("price", gorm.Expr("price * 0.9"))
		http.Redirect(w, r, "/admin/Product", http.StatusSeeOther)
	})
```

## Batch Actions

Batch actions run against selected record IDs.

```go
adm.Register(Product{}).
	AddBatchAction("batch_delete", "Delete Selected", func(res *admin.Resource, ids []string, w http.ResponseWriter, r *http.Request) {
		db.Where("id IN ?", ids).Delete(&Product{})
		http.Redirect(w, r, "/admin/Product", http.StatusSeeOther)
	})
```

## Dashboard Charts

Charts appear on the dashboard and return labels plus numeric values.

```go
adm.AddChart("Users by Role", "pie", func(db *gorm.DB) ([]string, []float64) {
	var rows []struct {
		Role  string
		Count int64
	}

	db.Model(&User{}).Select("role, count(*) as count").Group("role").Scan(&rows)

	labels := make([]string, 0, len(rows))
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		labels = append(labels, row.Role)
		values = append(values, float64(row.Count))
	}

	return labels, values
})
```

## Custom Pages

Custom pages can render arbitrary HTML through a standard `http.HandlerFunc`.

```go
adm.AddPage("SystemStatus", "Administration", func(w http.ResponseWriter, r *http.Request) {
	content := template.HTML(`<p>Server status: online</p>`)
	view.RenderCustomPage(adm, w, r, "System Status", content)
})
```

## Configuration

Load optional YAML configuration from `admin.yml`.

```go
conf, err := admin.LoadConfig("admin.yml")
if err == nil && conf != nil {
	adm.SetConfig(conf)
}
```

Example:

```yaml
site_title: "My Admin"
default_per_page: 10
```
