# packages
The install command is not a package manager, it merely leverages your system’s package manager,
whether it be `apt`, `homebrew`, `dnf`, or `paman`.

All packages available to install are located in the [packages.toml](packages.toml) file. Packages
are grouped into lists to be installed together. In the TOML file, package groups are initialized as
an array of tables (`[[package_group]]`). Each package group should have a name and description.
Packages should be created as tables, with a name in the style `[package_group.package_name]`. Each
package should have a `name` key to display when installing the package. Other possible keys for
each package are `apt_name`, `brew_name`, `brew_cask_name`, `dnf_name`, and `pacman_name`.

## List of packages
### Core
- [bat](https://github.com/sharkdp/bat) — A prettier version of `cat`
- [curl](https://curl.se/) — A command-line tool for getting or sending data
- [exa](https://github.com/ogham/exa) — A modern replacement for `ls`
- [fd](https://github.com/sharkdp/fd) — A simple, fast and user-friendly alternative to `find`
- [fzf](https://github.com/junegunn/fzf) — A command-line fuzzy finder
- [git](https://git-scm.com/) — A distributed version control system
- [gnupg](https://gnupg.org/) — A complete and free implementation of the OpenPGP standard
- [node](https://nodejs.org/) — An open-source, cross-platform JavaScript runtime environment
- [openssl](https://github.com/openssl/openssl) — TLS/SSL and crypto library
- [ripgrep](https://github.com/BurntSushi/ripgrep) — A modern replacement for `grep`
- [tmux](https://github.com/tmux/tmux) — A terminal multiplexer
- [vim](https://github.com/vim/vim) — My favorite editor, customizable, extensible, and blazing fast
- [zsh](https://www.zsh.org/) — A better shell and scripting language

### Design
- [FFmpeg](https://ffmpeg.org/) — A complete, cross-platform solution to record, convert and stream
  audio and video

### GUI Core
- [Clean My Mac](https://cleanmymac.com/) — A simple Mac cleaner, speed booster, and health guard
- [Dropbox](https://dropbox.com/) — Nice and easy file hosting and sharing
- [Docker Desktop](https://docker.com/) — A platform for building, deploying, and managing
  containerized applications
- [IINA](https://iina.io/) — The modern media player for macOS
- [iTerm2](https://iterm2.com/) — A replacement for the default macOS Terminal
- [Little Snitch](https://www.obdev.at/products/littlesnitch/index.html) — A host-based application
  firewall for macOS
- [Mullvad VPN](https://mullvad.net/en/) — A good, reliable, and open source VPN
- [Parallels Desktop](https://www.parallels.com/products/desktop/) — Seemless Windows/Linux parallel
  runtime for macOS
- [Raycast](https://www.raycast.com/) — A blazingly fast, totally extendable launcher
- [Zoom](https://zoom.us/) — A necessary (for-now) video-conferencing service

### GUI Design
- [Adobe Creative Cloud](https://www.adobe.com/creativecloud.html) — Terrible software, but I
  haven't found anything better
- [Glyphs](https://glyphsapp.com/) — A great macOS font editor
- [Rhino 3D](https://www.rhino3d.com/) — Graphics and computer-aided design software
- [Cycling '74 Max/MSP](https://cycling74.com/products/max) — A visual programming language for
  music and multimedia
