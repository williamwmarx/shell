# shell

My personal dotfiles and install script.

## Install 🚀

Install is easy, just run this one liner for a nice TUI that shows your options.

```bash
sh <(curl %INSTALL_URL%)
```

> See [INSTALL.md](INSTALL.md) for more options

## Dotfiles 🧩

%DOTFILES%

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

%PACKAGES%
