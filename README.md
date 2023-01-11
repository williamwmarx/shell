# shell

My dotfiles and install script.

## Install 🚀

Install is easy, just run this one-liner for a nice TUI that shows your options.

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
initialized as tables (`[package_group.packages.package_name]`). Each package should have a `name`
key to display when installing the package, as well as `description` and `url` keys to assist in
building this README. Other possible keys for each package are `apt_name`, `brew_name`,
`brew_cask_name`, `dnf_name`, `pacman_name`, and `install_command`.

### Core

Core packages I use daily

- [bat](https://github.com/sharkdp/bat) - A cat clone with wings.
- [cURL](https://curl.se/) - A command line tool for transferring data with URL syntax.
- [exa](https://the.exa.website/) - A modern replacement for ls.
- [fd](https://github.com/sharkdp/fd) - A simple, fast and user-friendly alternative to find.
- [fzf](https://github.com/junegunn/fzf) - A command-line fuzzy finder.
- [git](https://git-scm.com/) - A fast, scalable, distributed revision control system.
- [GnuPG](https://gnupg.org/) - A complete and free implementation of the OpenPGP standard.
- [Node](https://nodejs.org/) - A JavaScript runtime built on Chrome's V8 JavaScript engine.
- [Oh My Zsh](https://ohmyz.sh/) - A delightful community-driven framework for managing your zsh configuration.
- [OpenSSL](https://www.openssl.org/) - A robust, commercial-grade, and full-featured toolkit for the Transport Layer Security (TLS) and Secure Sockets Layer (SSL) protocols.
- [ripgrep](https://github.com/BurntSushi/ripgrep) - A line-oriented search tool that recursively searches your current directory for a regex pattern.
- [tmux](https://github.com/tmux/tmux) - A terminal multiplexer.
- [Vim](https://www.vim.org/) - A highly configurable text editor built to enable efficient text editing.
- [Zsh](https://www.zsh.org/) - A shell designed for interactive use, although it is also a powerful scripting language.

### Design

Packages for visual and sound design.

- [FFmpeg](https://ffmpeg.org/) - A complete, cross-platform solution to record, convert and stream audio and video.

### GUI Core

GUI apps I use daily.

- [Clean My Mac](https://macpaw.com/cleanmymac) - A macOS app to clean up your Mac.
- [Docker Desktop](https://www.docker.com/products/docker-desktop) - A desktop app for MacOS and Windows machines for the building and sharing of containerized applications and microservices.
- [Dropbox](https://www.dropbox.com/) - A file hosting service that offers cloud storage, file synchronization, personal cloud.
- [IINA](https://iina.io/) - The modern video player for macOS.
- [iTerm2](https://iterm2.com/) - A terminal emulator for macOS that does amazing things.
- [Little Snitch](https://www.obdev.at/products/littlesnitch/index.html) - A firewall application for macOS that monitors outgoing network connections and allows or denies them.
- [Mullvad VPN](https://mullvad.net/en/) - A VPN service that helps keep your online activity, identity, and location private.
- [Parallels Desktop](https://www.parallels.com/products/desktop/) - A software for running Windows, Linux, or any other operating system on a Mac without rebooting.
- [Raycast](https://raycast.com/) - A smart command line productivity tool for macOS.
- [Zoom](https://zoom.us/) - A video conferencing, online chat, and web conferencing software.

### GUI Design

GUI apps for visual and sound design.

- [Adobe Creative Cloud](https://www.adobe.com/creativecloud.html) - A collection of desktop and mobile apps and services for photography, design, video, web, UX and more.
- [Cycling ‘74 Max/MSP](https://cycling74.com/products/max/) - A visual programming language for music, audio, and multimedia.
- [Glyphs](https://glyphsapp.com/) - A powerful font editor for macOS.
- [Rhino 3D](https://www.rhino3d.com/) - A 3D modeling software for Windows, macOS, and Linux.
