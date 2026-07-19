# Publishing

Go packages are published by pushing a public Git repository and tagging semantic versions.

## Repository Path

The module path in `go.mod` must match the public repository URL.

```go
module github.com/go-packs/go-admin
```

For this module path, the GitHub repository should be:

```text
https://github.com/go-packs/go-admin
```

## Release Checklist

- Keep `README.md`, `USAGE.md`, `examples/`, tests, CI, and `LICENSE`.
- Do not commit generated databases, uploads, build output, local caches, secrets, or `.env` files.
- Run `go test ./...`.
- Commit the release-ready changes.
- Tag a version such as `v0.1.0`.

## Publish A Version

```bash
git add .
git commit -m "Prepare open source release"
git push origin main

git tag v0.1.0
git push origin v0.1.0
```

Users can install the tagged release with:

```bash
go get github.com/go-packs/go-admin@v0.1.0
```

## Package Documentation

Go API documentation appears on pkg.go.dev after the package is public:

```text
https://pkg.go.dev/github.com/go-packs/go-admin
```

## Showcase Documentation

This docs site is published with GitHub Pages. After the workflow runs, the site URL is:

```text
https://go-packs.github.io/go-admin/
```

If the repository lives under a different GitHub user or organization, update `site_url`, `repo_url`, and `repo_name` in `mkdocs.yml`.
