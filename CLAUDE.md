# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based dotfiles management and system configuration tool that provides both interactive (TUI) and non-interactive installation of dotfiles and packages across different Unix-like systems.

## Architecture

### Core Components

- **Main Entry Point**: `main.go` → `cmd.Execute()`
- **CLI Framework**: Built with Cobra for command handling and Bubble Tea for interactive TUI
- **Configuration**: TOML-based configuration in `config.toml` defines:
  - Dotfile sync targets (mapping repo paths to local paths)
  - Package definitions with package manager mappings (apt, brew, dnf, pacman)
  - Installer definitions for different components (vim, tmux, zsh)

### Key Modules

- `cmd/root.go`: CLI entry point and flag parsing
- `cmd/tui.go`: Interactive terminal UI implementation
- `cmd/installers.go`: Installation logic for packages and dotfiles
- `cmd/helpers.go`: Utility functions for downloads, command execution, and configuration parsing
- `assembly/main.go`: Assembly functions for preparing and building the application

## Build and Development Commands

```bash
# Build the binary
go build

# Run the binary
./shell --help

# Test the TUI
./shell

# Build for release (cross-platform)
make build  # outputs to releases/

# Test in Docker (Linux environment)
make test -- --zsh
make test -- --vim --tmux

# Regenerate README.md and INSTALL.md from templates
go run assembly/main.go

# Format code before committing
gofmt -s -w .
```

## Installation Flow

1. User runs `sh <(curl https://marx.sh)` which downloads and executes `install.sh` (via Cloudflare Worker)
2. `install.sh` fetches the pre-built binary for the user's platform from GitHub Releases
3. The binary presents either:
   - Interactive TUI (default)
   - Non-interactive installation with flags (`--vim`, `--tmux`, `--zsh`, `--full`)
4. Selected packages are installed using the system's package manager
5. Dotfiles are synced from the repo to local paths defined in `config.toml`

## Package Management

The tool supports multiple package managers:
- **macOS**: Homebrew (brew, brew cask)
- **Linux**: apt, dnf, pacman
- **Fallback**: Custom install commands

Package definitions are in `packages.toml` with mappings for each package manager.

## Repository URL Configuration

The repository owner and name are configured in `cmd/helpers.go`:
- Owner: "williamwmarx"
- Repo: "shell"
- These are hard-coded to fix installation issues when running from non-repository directories