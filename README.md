# shell
My personal dotfiles and install script.

## Install 🚀
Install is easy, just run this one liner for a nice TUI that shows your options.
```bash
sh <(curl https://marx.sh)
```
> *See [INSTALL.md](INSTALL.md) for more options*

## Components 🧩
### Config Files
- Git
  - [.gitconfig](git/gitconfig) — Set user info, gpg signing, and “main” as default branch
- GnuPG
  - [.gpg.conf](gnupg/gpg.conf) — macOS pinentry program
  - [.gpg-agent.conf](gnupg/gpg-agent.conf) — macOS pinentry program
- Tmux
  - [.tmux.conf](tmux/tmux.conf) — Clean UI, useful info only, Vim-like keybindings
- Vim
  - [.vimrc](vim/vimrc) — My day-to-day Vim config, batteries included
  - [Template files](vim/templates) — Skeleton templates when creating a new file
- Zsh
  - [.zshrc](zsh/zshrc) — My go-to shell setup, leverages [Oh-My-Zsh](https://ohmyz.sh/)
  - [.aliases](zsh/aliases) — Useful aliases, built to work across all UNIX systems
  - [.functions](zsh/functions) — Useful functions, built to work across all UNIX systems
  - [t3.zsh-theme](zsh/t3.zsh-theme) — A nice, clean theme with everything you need, and nothing more

### Packages
- Core
  - [bat](https://github.com/sharkdp/bat) — A prettier version of `cat`
  - [curl](https://curl.se/) — A command-line tool for getting or sending data
  - [exa](https://github.com/ogham/exa) — A modern replacement for `ls`
  - [fd](https://github.com/sharkdp/fd) — A simple, fast and user-friendly alternative to `find`
  - [fzf](https://github.com/junegunn/fzf) — A command-line fuzzy finder
  - [git](https://git-scm.com/) — A distributed version control system
  - [gnupg](https://gnupg.org/) — A complete and free implementation of the OpenPGP standard
  - [openssl](https://github.com/openssl/openssl) — TLS/SSL and crypto library
  - [ripgrep](https://github.com/BurntSushi/ripgrep) — A modern replacement for `grep`
  - [tmux](https://github.com/tmux/tmux) — A terminal multiplexer
  - [vim](https://github.com/vim/vim) — My favorite editor, customizable, extensible, and blazing fast
  - [zsh](https://www.zsh.org/) — A better shell and scripting language
- Design
  - [FFmpeg](https://ffmpeg.org/) — A complete, cross-platform solution to record, convert and stream audio and video
- GUI Core
  - [Clean My Mac](https://cleanmymac.com/) — A simple Mac cleaner, speed booster, and health guard
  - [Dropbox](https://dropbox.com/) — Nice and easy file hosting and sharing
  - [Docker](https://docker.com/) — A platform for building, deploying, and managing containerized applications
  - [IINA](https://iina.io/) — The modern media player for macOS
  - [iTerm2](https://iterm2.com/) — A replacement for the default macOS Terminal
  - [Little Snitch](https://www.obdev.at/products/littlesnitch/index.html) — A host-based application firewall for macOS
  - [Mullvad VPN](https://mullvad.net/en/) — A good, reliable, and open source VPN
  - [Parallels Desktop](https://www.parallels.com/products/desktop/) — Seemless Windows/Linux parallel runtime for macOS
  - [Raycast](https://www.raycast.com/) — A blazingly fast, totally extendable launcher
  - [Zoom](https://zoom.us/) — A necessary (for-now) video-conferencing service
- GUI Design
  - [Adobe Creative Cloud](https://www.adobe.com/creativecloud.html) — Terrible software, but I haven't found anything better
  - [Glyphs](https://glyphsapp.com/) — A great macOS font editor
  - [Rhino 3D](https://www.rhino3d.com/) — Graphics and computer-aided design software
  - [Cycling '74 Max/MSP](https://cycling74.com/products/max) — A visual programming language for music and multimedia
  
### Other
- macOS
  - [macOS](macOS/macOS) — macOS defaults to set on new machine
- Personal
  - [.plan](personal/plan) — A nostalgic resurrection of the [Carmack-esque](https://garbagecollected.org/2017/10/24/the-carmack-plan/) .plan file
- Raycast
  - [clear-format.sh](raycast/clear-format.sh) — Strip formatting (while keeping whitespace) from clipboard
  - [expand-url.sh](raycast/expand-url.sh) — Follow URL redirects, returning final location
