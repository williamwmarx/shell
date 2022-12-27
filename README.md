# shell
My personal dotfiles and install script.

## Install 🚀
Install is easy, just run this one liner for a nice TUI that shows your options.
```bash
sh <(curl https://marx.sh)
```
> See [INSTALL.md](INSTALL.md) for more options

## Packages 📦
See the [packages](packages) README for more info

## Dotfiles 🧩
### Git
- [.gitconfig](git/gitconfig) — Git configuration

### GnuPG
Set up pinentry program
- [gpg.conf](gnupg/gpg.conf) — GnuPG configuration
- [gpg-agent.conf](gnupg/gpg-agent.conf) — GnuPG agent configuration

### Raycast
Raycast shell scripts
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
- [.zshrc](zsh/zshrc) — My go-to shell setup, leveraging Oh My Zsh (if available)
- [.aliases](zsh/aliases) — Useful aliases, built to work across all UNIX systems
- [.functions](zsh/functions) — Useful functions, built to work across all UNIX systems
- [t3.zsh-theme](zsh/t3.zsh-theme) — A nice, clean theme with everything you need, and nothing more