# lenslocked

Follow-along Go web development project and course material. This repository contains the code examples, exercises, and supporting materials used in the "Lenslocked" Go web development tutorial.

## Overview

Lenslocked is a sample web application used to teach building production-quality web applications in Go. It demonstrates:

- Project layout for Go web applications
- Dependency management with Go modules (go.mod)
- HTTP routing, handlers, and middleware
- Templates and static assets handling
- Database access patterns and migrations
- Testing, tooling, and best practices for maintainable code

This repo is intended to be used as both a learning resource and a starting point for small Go web projects.

## Repository structure (typical)

Note: your local copy of this repository may vary. Common folders you can expect in a Go web app:

- cmd/         - application entrypoints
- internal/    - private application code
- pkg/         - reusable packages (optional)
- web/         - templates, static assets, front-end resources
- migrations/  - database migration files
- go.mod       - module and dependency management

If a particular file or folder above is not present, consult the code to find the actual layout.

## Prerequisites

- Go 1.18+ installed (use the current stable release when possible)
- Git for cloning the repository
- (Optional) A database server if the app uses one locally (Postgres, MySQL, SQLite, etc.)

## Getting started (local development)

1. Clone the repository:

   git clone https://github.com/raminderis/lenslocked.git
   cd lenslocked

2. Install dependencies:

   go mod download

3. Build the project:

   go build ./...

4. Run the application:

- If the repository has a single main package at the repository root, you can run:

      go run .

- If the project exposes one or more entrypoint binaries under cmd/, run the desired command, for example:

      go run ./cmd/<binary-name>

Replace `<binary-name>` with the actual folder name under cmd/.

5. Running tests:

   go test ./...

6. Formatting and vetting code:

   go fmt ./...
   go vet ./...

## Development workflow

- main: stable branch containing release-ready code and documentation.
- dev: active development branch (feature branches and PRs should be opened against dev).

Typical flow:

1. Create a topic branch from dev: `git checkout -b feat/your-feature dev`
2. Implement changes and add tests.
3. Run `go test ./...`, `gofmt`, and `go vet` locally.
4. Push your branch and open a pull request targeting `dev`.
5. After review and testing, merge into `dev`. Periodically `dev` is merged into `main` for releases.

If your team uses a different branching strategy, follow the repository CONTRIBUTING.md or team conventions.

## Testing and CI

- Use `go test ./...` to run unit tests.
- Add integration tests where appropriate and document how to run them (e.g., requiring a running DB).
- CI pipelines (GitHub Actions or other) should run `go test`, `go vet`, and `gofmt` checks.

## Contributing

Contributions are welcome. Please follow these general guidelines:

- Open an issue to discuss large or unclear changes before implementing.
- Fork and create a feature branch from `dev` for code changes.
- Include tests for new features or bug fixes.
- Keep changes focused and provide a clear, descriptive PR title and description.
- Run `gofmt` before submitting a PR.

For more specific rules (code style, commit message format, review process), add a CONTRIBUTING.md to the repository.

## License

Check the repository root for a LICENSE file. If no license is present, the project has no explicit open-source license — consider adding one (MIT / Apache 2.0 are common choices) to make reuse and contributions straightforward.

## Contact / Author

Maintainer: raminderis (GitHub)

For questions, issues, or requests, open an issue in this repository or contact the maintainer via their GitHub profile.

---

Thank you for checking out lenslocked — good luck learning and building with Go!
