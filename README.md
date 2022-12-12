# shell
My personal dotfiles and install script.

## Install 🚀
Install is easy, just run this one liner for a nice TUI that shows your options.
```bash
sh <(curl https://marx.sh)
```
> See [INSTALL.md](INSTALL.md) for more options

## Components 🧩
### Git
- [.gitconfig](git/gitconfig) — Set user info, gpg signing, and “main” as default branch

### GnuPG
- [.gpg.conf](gnupg/gpg.conf) — macOS pinentry program
- [.gpg-agent.conf](gnupg/gpg-agent.conf) — macOS pinentry program

### macOS
- [macOS](macOS/macOS) — macOS defaults to set on new machine

### Packages
- See the [packages](packages) README for more info

### Personal
- [.plan](personal/plan) — A nostalgic resurrection of the [Carmack-esque](https://garbagecollected.org/2017/10/24/the-carmack-plan/) .plan file

### Raycast
- [clear-format.sh](raycast/clear-format.sh) — Strip formatting (while keeping whitespace) from clipboard
- [expand-url.sh](raycast/expand-url.sh) — Follow URL redirects, returning final location

### skhd
- [.skhdrc](skhd/skhdrc) — Configure the hotkey daemon I use (mainly to manage Yabai)

### Tmux
- [.tmux.conf](tmux/tmux.conf) — Clean UI, useful info only, Vim-like keybindings

### Vim
- [.vimrc](vim/vimrc) — My day-to-day Vim config, batteries included
- [Template files](vim/templates) — Skeleton templates when creating a new file

### Yabai
- [.yabairc](yabai/yabairc) — Manage Yabai, the tiling window manager I use

### Zsh
- [.zshrc](zsh/zshrc) — My go-to shell setup, leverages [Oh-My-Zsh](https://ohmyz.sh/)
- [.aliases](zsh/aliases) — Useful aliases, built to work across all UNIX systems
- [.functions](zsh/functions) — Useful functions, built to work across all UNIX systems
- [t3.zsh-theme](zsh/t3.zsh-theme) — A nice, clean theme with everything you need, and nothing more
