<section class="hero">
  <h1>Build a clean admin panel for your Go models.</h1>
  <p>Go Admin uses GORM metadata and a small registration API to generate CRUD screens, dashboards, authentication, audit logs, filters, actions, and exports.</p>
  <div class="hero-actions">
    <a href="quick-start/">Start building</a>
    <a href="examples/">Run the demo</a>
    <a href="https://pkg.go.dev/github.com/go-packs/go-admin">API docs</a>
  </div>
</section>

## Why Go Admin

<div class="feature-grid">
  <div class="feature-card">
    <h3>Fast setup</h3>
    <p>Register a GORM model, expose fields, and mount the admin router.</p>
  </div>
  <div class="feature-card">
    <h3>Built-in auth</h3>
    <p>Session-based login with admin users, roles, permissions, and bcrypt passwords.</p>
  </div>
  <div class="feature-card">
    <h3>Useful workflows</h3>
    <p>Batch actions, CSV export, search filters, custom pages, and audit logging.</p>
  </div>
  <div class="feature-card">
    <h3>Portable UI</h3>
    <p>Templates and styles are embedded into the Go binary with `go:embed`.</p>
  </div>
</div>

## Install

```bash
go get github.com/go-packs/go-admin
```

## Minimal Example

```go
adm := admin.NewRegistry(db)

adm.Register(Product{}).
	SetGroup("Inventory").
	RegisterField("ID", "ID", true).
	RegisterField("Name", "Product Name", false).
	RegisterField("Price", "Price", false)

http.Handle("/admin/", server.NewRouter(adm))
log.Fatal(http.ListenAndServe(":8080", nil))
```

## Documentation Map

- [Quick Start](quick-start.md): create a small app and mount the admin panel.
- [Usage Guide](usage.md): fields, associations, actions, batch actions, charts, and pages.
- [Examples](examples.md): run the included demo app locally.
- [Publishing](publishing.md): prepare and publish the package as open source.
