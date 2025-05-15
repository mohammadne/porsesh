# Porsesh

![Go Version](https://img.shields.io/badge/Golang-1.24-66ADD8?style=for-the-badge&logo=go)
![App Version](https://img.shields.io/github/v/tag/mohammadne/porsesh?sort=semver&style=for-the-badge&logo=github)
![Repo Size](https://img.shields.io/github/repo-size/mohammadne/porsesh?logo=github&style=for-the-badge)
![Coverage](https://img.shields.io/codecov/c/github/mohammadne/porsesh?logo=codecov&style=for-the-badge)

> Porsesh is a Persian word that means "question" — a fitting name for a compact yet powerful Polling platform.

The tiny polling platform

## Usage

### Local

```sh
# run dependencies
cd hacks/compose && podman compose -f compose.local.yml up -d

# run migration files
go run cmd/migration/main.go --direction=up

# run server
go run cmd/server/main.go
```
