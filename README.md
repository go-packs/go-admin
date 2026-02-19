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
- 🚀 **Portable**: Everything (HTML/CSS/JS) is bundled into your binary using `go:embed`.

## Installation

```bash
go get github.com/ajeet-kumar1087/go-admin
```

## Usage Guide

### 1. Basic Setup

Initialize the registry with a GORM database connection and handle the `/admin/` route.

```go
db, _ := gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
adm := admin.NewRegistry(db)

// Start Server
http.Handle("/admin/", adm)
http.ListenAndServe(":8080", nil)
```

### 2. Registering Resources

Register a struct to create a CRUD interface. You can specify which fields are shown and which are read-only.

```go
adm.Register(Product{}).
    SetGroup("Inventory").
    RegisterField("ID", "ID", true).
    RegisterField("Name", "Product Name", false).
    RegisterField("Price", "Price", false).
    SetIndexFields("Name", "Price") // Only show these on the list page
```

### 3. Adding Scopes (Tabs)

Scopes provide predefined filters at the top of your resource lists.

```go
adm.Register(Product{}).
    AddScope("expensive", "Expensive Items", func(db *gorm.DB) *gorm.DB {
        return db.Where("price > ?", 1000)
    })
```

### 4. Custom Actions

#### Member Actions (Single Record)
Shown on the detail (Show) page of a specific record.

```go
resource.AddMemberAction("activate", "Activate User", func(res *admin.Resource, w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    // Your logic here
    http.Redirect(w, r, "/admin/User", http.StatusSeeOther)
})
```

#### Batch Actions (Multiple Records)
Allows selecting rows via checkboxes and applying an operation.

```go
resource.AddBatchAction("bulk_delete", "Delete Selected", func(res *admin.Resource, ids []string, w http.ResponseWriter, r *http.Request) {
    db.Where("id IN ?", ids).Delete(&Product{})
    http.Redirect(w, r, "/admin/Product", http.StatusSeeOther)
})
```

### 5. Associations

Go Admin handles relationships automatically.

```go
// HasMany: Product has many Specifications
productRes.HasMany("Specs", "Technical Specs", "ProductInfo", "ProductID")

// BelongsTo: ProductInfo belongs to a Product
infoRes.BelongsTo("ProductID", "Parent Product", "Product", "ID")
```

### 6. Dashboard Charts

Add visual widgets to your home page.

```go
adm.AddChart("Sales by Category", "pie", func(db *gorm.DB) ([]string, []float64) {
    // Return labels and values from DB
    return []string{"Electronics", "Books"}, []float64{500, 200}
})
```

## Security Note

By default, Go Admin creates an `AdminUser` table. Ensure you seed a default user to log in:

```go
adminUser := &admin.AdminUser{Email: "admin@example.com", Role: "admin"}
adminUser.SetPassword("securepassword")
db.Create(adminUser)
```

## License

This project is licensed under the MIT License.
