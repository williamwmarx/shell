# tmux
## Configuration π§©
There're too many choices to describe, but the [.tmux.conf](tmux.conf) file itself is well-documented.

### Keybindings
Keybindings are Vi-style
- Split panes
    - `Ctrl-b + s` βΒ Split window horizontally 
    - `Ctrl-b + v` βΒ Split window vertically 
- Pane navigation
    - `Ctrl-b + h` β Move left
    - `Ctrl-b + j` β Move down
    - `Ctrl-b + k` β Move up
    - `Ctrl-b + l` β Move right
- Reize panes
    - `Ctrl-b + <` β Decrease horizontally split pane size
    - `Ctrl-b + >` β Increase horizontally split pane size
    - `Ctrl-b + -` β Decrease vertically split pane size
    - `Ctrl-b + +` β Increase vertically split pane size
- Copy/Paste
    - `Ctrl-b + [` β Enter copy/paste mode (navigate with Vi-style keys)
    - `Ctrl-b + v` β Start selecting text
    - `Ctrl-b + y` β Copy selection
    - `Ctrl-b + P` β Paste selection
- Miscellaneous
    - `Ctrl-b + x` β List current tmux sessions

## Theme π¨
Like the Vim and Zsh themes, only what's necessary is shown. The theme itself is dark, but when a pane is inactive, the background lightens slightly.

The bottom bar shows the session name, all open windows, the current username and hostname, and the date and time. The name of each window follows the format `window_number:process_name`. If a window is active, it's format is `window_number[pane_number]:process_name`.

<img src="../assets/TmuxThemePreview.png" alt="Tmux Preview" width="80%"/>
