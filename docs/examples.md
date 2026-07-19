# Examples

The repository includes a runnable demo app under `examples/`.

## Run The Demo

```bash
go run ./examples
```

The demo creates a local SQLite database at:

```text
examples/demo.db
```

Open the admin panel:

```text
http://localhost:8080/admin
```

Sign in with:

```text
Email: admin@example.com
Password: password123
```

## What The Demo Shows

- Admin users, roles, permissions, sessions, and audit logs.
- Product and user resources.
- Select fields, number fields, image fields, and decorators.
- Has-many and belongs-to associations.
- Collection actions, member actions, and batch actions.
- Dashboard charts and custom pages.

## Reset Demo Data

Stop the server, delete the generated database, then run the example again.

```bash
rm examples/demo.db
go run ./examples
```

The database file is intentionally ignored by Git. Users should generate local data when they run examples.
