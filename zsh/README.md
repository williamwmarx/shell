# Zsh
## Configuration 🧩
### zshrc
- Install (if necessary) [Oh-My-Zsh](https://ohmyz.sh/) and the two plugins we use, enabling [autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) and [syntax highlighting](https://github.com/zsh-users/zsh-syntax-highlighting)
- Export the PATH variable, exposing the binaries we installed to the command line
- Configure [fzf](https://github.com/junegunn/fzf)
- Vimify everything
  - Vi mode to navigate command line
    - `v` in escape mode launches Vim with current string
  - Use [Vim](https://github.com/vim/vim) for manpages
- Load all [aliases](aliases) and [functions](functions)
- Launch [tmux](https://github.com/tmux/tmux) on start
  - If no session present, starts new session with name `main`
  - If session present, starts new session with name in format `alt[n_sessions - 1]` (e.g. `alt0`, `alt1`, etc.)

## Theme 🎨
There’s not much to my default Zsh theme — just a simple left prompt showing the current path and a right prompt showing git status if available. For the left prompt, if the user is `root`, the leader, `>`, changes from bold white to bold red.

The below image additionally shows off the usefullness of [Zsh syntax highlighting](https://github.com/zsh-users/zsh-syntax-highlighting).

![Zsh Theme Preview](../assets/ZshThemePreview.png)
