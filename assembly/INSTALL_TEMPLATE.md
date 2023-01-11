# Install

The recommended way to install anything is to run this one-liner, which loads a nice TUI showing
your options.

```bash
sh <(curl %INSTALL_URL%)
```

## Non-interactive install options

### Install everything

```bash
sh <(curl %INSTALL_URL%) --full
```

### Partial install

Any of the flags below can be combined.

%PARTIAL_INSTALL%

### Temporary install

Sometimes, you only need your dotfiles temporarily. For example, say you're editing some code on a
friend's machine. You could slowly go through it with their editor, or you could load up your vim
config and fly through their code. This is where the `--tmp` flag comes in. You can use the `--tmp`
flag with %TMP_FLAGS%. It will install the packages, download necessary dotfiles into the
`%TMP_DIR%` directory, and add the shell script `%TMP_DIR%/uninstall.sh` which will uninstall any
packages you installed and remove the `%TMP_DIR%` directory. Temporary install will look for the
“vanilla” versions of synced dotfiles, where possible.
