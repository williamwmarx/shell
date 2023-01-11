# shell

My personal dotfiles and install script.

## Install 🚀

Install is easy, just run this one liner for a nice TUI that shows your options.

```bash
sh <(curl https://marx.sh)
```

> See [INSTALL.md](INSTALL.md) for more options

## Dotfiles 🧩

### Git

- [.gitconfig](git/gitconfig) — Git configuration

### GnuPG

- [gpg-agent.conf](gnupg/gpg-agent.conf) — GnuPG agent configuration
- [gpg.conf](gnupg/gpg.conf) — GnuPG configuration

### Raycast

- [clear-format.sh](raycast/clear-format.sh) — Strip formatting (while preserving whitespace) from clipboard

### skhd

- [.skhdrc](skhd/skhdrc) — Configure the hotkey daemon I use (mainly to manage Yabai)

### tmux

- [.tmux.conf](tmux/tmux.conf) — Clean UI, useful info only, Vim-like keybindings

### Vim

- [.vimrc](vim/vimrc) — Vim configuration

### yabai

- [.yabairc](yabai/yabairc) — Manage yabai, the tiling window manager I use

### Zsh

- [.aliases](zsh/aliases) — Useful aliases, built to work across all UNIX systems
- [.functions](zsh/functions) — Useful functions, built to work across all UNIX systems
- [.zshrc](zsh/zshrc) — My go-to shell setup, leveraging Oh My Zsh (if available)
- [t3.zsh-theme](zsh/t3.zsh-theme) — A nice, clean theme with everything you need, and nothing more

## Packages 📦

The install command is not a package manager, it merely leverages your system’s package manager,
whether it be `apt`, `homebrew`, `dnf`, or `pacman`. It can also run install scripts for packages
not in the above indices.

All packages available to install are located in the [packages.toml](packages.toml) file. Packages
are grouped into lists to be installed together. In the TOML file, packages and package groups are
initialized as tables (`[package_group.package]`). Each package should have a `name` key to display
when installing the package, as well as `description` and `url` keys to assist in building this
README. Other possible keys for each package are `apt_name`, `brew_name`, `brew_cask_name`,
`dnf_name`, `pacman_name`, and `install_command`.

### Core

Core packages I use daily

### Design



### GUI Core



### GUI Design
