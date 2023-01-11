# Install

The recommended way to install anything is to run this one liner, which loads a nice TUI showing
your options.

```bash
sh <(curl https://marx.sh)
```

## Non-interactive install options

### Install everything

```bash
sh <(curl https://marx.sh) --full
```

### Partial install

Any of the flags below can be combined.
#### tmux
Install [tmux](https://github.com/tmux/tmux) and [.tmux.conf](tmux/tmux.conf)
```bash
sh <(curl https://marx.sh) --tmux
```

#### vim
Install [Vim](https://github.com/vim/vim), my [.vimrc](vim/vimrc), and [skeleton files](vim/templates)
```bash
sh <(curl https://marx.sh) --vim
```

#### zsh
Install [Zsh](https://www.zsh.org/), [Oh-My-Zsh](https://ohmyz.sh/), [.zshrc](zsh/zshrc), [.aliases](zsh/aliases), [.functions](zsh/functions), and [Zsh theme](zsh/t3.zsh-theme)
```bash
sh <(curl https://marx.sh) --zsh
```

### Temporary install
Sometimes, you only need your dotfiles temporarily. For example, say you're editing some code on a friend's machine. You could slowly go through it with their editor, or you could load up your vim config and fly through their code. This is where the `--tmp` flag comes in. You can use the `--tmp` flag with `--tmux`, `--vim`, or `--zsh`. It will install the packages, download necessary dotfiles into the `~/.shell.tmp` directory, and add the shell script `~/.shell.tmp/uninstall.sh` which will uninstall any packages you installed and remove the `~/.shell.tmp` directory. Temporary install will look for the “vanilla” versions of synced dotfiles, where possible.