# shell

> **⚠️ This repository is archived and no longer maintained.**
> I've moved my dotfiles to a private repository using [GNU Stow](https://www.gnu.org/software/stow/) for simpler symlink management.

A Go-based dotfiles management tool featuring an interactive TUI for installing dotfiles and packages across Unix-like systems.

## What This Was

This project was a learning exercise to:
- Get familiar with Go conventions and patterns
- Build a CLI with [Cobra](https://github.com/spf13/cobra) and an interactive TUI with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Solve the cross-platform dotfiles management problem I had at the time

### Features

- Interactive terminal UI for selecting what to install
- Non-interactive mode with flags (`--vim`, `--tmux`, `--zsh`, `--full`)
- Cross-platform package manager support (apt, brew, dnf, pacman)
- TOML-based configuration for dotfile mappings and package definitions

## Dotfiles

The dotfiles remain in this repo for reference:

- **Git**: [gitconfig](git/gitconfig)
- **GnuPG**: [gpg.conf](gnupg/gpg.conf), [gpg-agent.conf](gnupg/gpg-agent.conf)
- **tmux**: [tmux.conf](tmux/tmux.conf)
- **Vim**: [vimrc](vim/vimrc)
- **Zsh**: [zshrc](zsh/zshrc), [aliases](zsh/aliases), [functions](zsh/functions), [t3.zsh-theme](zsh/t3.zsh-theme)
- **yabai/skhd**: [yabairc](yabai/yabairc), [skhdrc](skhd/skhdrc)

## License

Apache 2.0
