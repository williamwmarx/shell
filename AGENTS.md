# Repository Guidelines

## Project Structure & Module Organization
- `cmd/`: Go CLI (Cobra + Bubble Tea) logic, installers, helpers, TUI.
- `main.go`: CLI entrypoint calling `cmd.Execute()`.
- `assembly/`: Generators for `README.md` and `INSTALL.md`; run to refresh docs.
- `assets/`: Images used in documentation.
- `config.toml`: Sync targets (dotfiles) and installer definitions.
- `packages.toml`: Package groups and per–package metadata.
- Dotfiles: `git/`, `gnupg/`, `raycast/`, `skhd/`, `tmux/`, `vim/`, `yabai/`, `zsh/`.
- CI: `.github/workflows/markdown.yml` (docs build) and `deploy.yml` (release build + S3 sync).

## Build, Test, and Development Commands
- Build releases: `make build` → binaries in `releases/` (Darwin/Linux arches).
- Clean: `make clean` → remove build artifacts.
- Run locally (TUI): `go run main.go` (use arrows/enter).
- Non‑interactive runs: `go run main.go --full` or `go run main.go --vim --tmp`.
- Docker test harness: `make test -- --zsh` or `make test -- --vim --tmp`.
- Regenerate docs: `go run assembly/main.go` (updates `README.md`, `INSTALL.md`).

## Coding Style & Naming Conventions
- Go: follow `gofmt`/idiomatic Go. Run `gofmt -s -w .` before committing.
- Packages lowercase; exported identifiers `CamelCase`; flags are long `--flag` names.
- Keep CLI messages concise, present tense. Shell snippets should be POSIX‐sh compatible.

## Testing Guidelines
- No unit tests yet. If adding tests, use Go’s `testing` with `*_test.go` and run `go test ./...`.
- Prefer validating flows via non‑interactive flags (examples above). Capture reproducible steps in PRs.

## Commit & Pull Request Guidelines
- Use clear, scoped commits; Conventional Commits are welcome (e.g., `feat:`, `fix:`, `chore:`).
- PRs should include: summary, rationale, manual test steps/commands, and linked issues.
- If you change `config.toml`, `packages.toml`, or `assembly/` templates, run `go run assembly/main.go` and commit regenerated markdown.
- Include screenshots when altering assets/UI previews. Avoid committing secrets or machine‑specific paths.

## Security & Configuration Tips
- The installer fetches files and may run package manager commands; review changes carefully.
- Prefer `--tmp` for ephemeral setups; it writes to `~/.@repo_name.tmp` and includes an uninstall script.
