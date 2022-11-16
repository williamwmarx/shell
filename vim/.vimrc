""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" Information
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"Maintained by William W. Marx (m94@marx.sh)
"🄯 Copyleft 2022, All Wrongs Reserved
"https://github.com/williamwmarx/.dotfiles

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" General
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
set nocompatible	"Don't be a luddite
filetype plugin indent on  "All filetype detection capabilities
set encoding=utf-8	"Use utf-8 as encoding
set fileformats=unix,dos,mac	"Set standard filetype to UNIX
if empty(glob('$HOME/.vim/swp'))	"Make .vim/swp directory if necessary
	silent !mkdir $HOME/.vim/swp
endif
set directory=$HOME/.vim/swp//	"Set swapfile directory
set nobackup	"Don't use backupfiles
set nowritebackup  "Don't write backupfiles
set updatetime=300	"Shorter update time
set signcolumn=yes	"Always show sign column
set hidden	"Don't force write before switching buffers
set history=500  "Keep X commands in history
set magic  "Turn magic on for RegEx
set bsdir=last	"Last accessed directory is default working directory
set incsearch  "Show search results while typing
set wildmenu	"Enhanced command line completion with tab
set wildmode=longest:list,full	"Complete to longest, then complete to first
set nojoinspaces	"If space at end of line, don't add another on join
set viewoptions-=options	"Don't let vim autochange directories with :mkview
set omnifunc=syntaxcomplete#Complete

" ---------------------------------- Spacing ----------------------------------
set tabstop=2  "Width of \t character
set shiftwidth=2 "Indentation width (>>, << chars)
set autoindent	"Copy indentaion from previous line
set textwidth=100  "100 chars per line. This isn't a VT100, darling
set nowrap	"Don't wrap text unless specified by filetype below
set backspace=indent,eol,start	"More powerful backspace

" ------------------------------- No Error Bells ------------------------------
set noerrorbells	"No sounds for errors
set novisualbell	"No screen flashes for errors
set t_vb=  "Ensure visual bell is off
set tm=500
if has("gui_macvim")
	autocmd GUIEnter * set vb t_vb=  "No sound on MacVim
endif

" --------------------------------- Aesthetics --------------------------------
set background=dark  "Ensure dark background
syntax on  "Synatx highlighting
set hlsearch	"Highlight search results
set number relativenumber  "Turn on relative line numbers
set list	"Show whitespace characters
set listchars=eol:$,tab:>-,space:·	"Add tab to list of displayed chars
set showcmd  "Show commands
set splitright	"Split vertical windows to the right
set splitbelow	"Split horizontal windows to the bottom
set colorcolumn=+1	"Put a vertical bar one char past texwidth
let g:netrw_banner=0	"Get rid of netrw banner
let g:netrw_liststyle=3  "Tree style directory listing
let g:netrw_winsize = 30	"Netrw takes up only 30% of the window

" --------------------------------- Functions ---------------------------------
func CleanRegs()
	let regs=split('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/-"', '\zs')
	for r in regs
		call setreg(r, [])
	endfor
endfunc

" -------------------------------- Status Line --------------------------------
let g:current_mode={'n': 'NORMAL', 'no': 'NORMAL·OPERATOR PENDING', 'v': 'VISUAL', 'V': 'V·LINE',
	\'^V': 'V·BLOCK', 's': 'SELECT', 'S': 'S·LINE', '^S': 'S·BLOCK', 'i': 'INSERT', 'R': 'REPLACE',
	\'Rv': 'V·REPLACE', 'c': 'COMMAND', 'cv': 'VIM EX', 'ce': 'EX', 'r': 'PROMPT', 'rm': 'MORE',
	\'r?': 'CONFIRM', '!': 'SHELL', 't': 'TERMINAL'}	"Mode abrv to word mappings

function! SLGitBranch()  "Get current git branch for status line
	let l:branch = system("git rev-parse --abbrev-ref HEAD 2>/dev/null | tr -d '\n'")
	return strlen(l:branch) > 0 ? '[' . l:branch . ']' : ''
endfunction

function! SLFileSize()	"Current buffer file size for status line
	let l:bytes = getfsize(expand('%:p'))
	if l:bytes >= 1024000
		return bytes / 1024000 .'MB'
	elseif l:bytes >= 1024
		return l:bytes / 1024 .'KB'
	elseif l:bytes >= 0
		return bytes .'B'
	else
		return '0B'
	endif
endfunction

set laststatus=2	"Always show status line
set statusline=  "Left status line
set statusline+=\ %{g:current_mode[mode()]}  "Operating mode
set statusline+=\ %f	"Buffer number and Filename
set statusline+=\ %{SLGitBranch()}	"Git branch
set statusline+=%=	"Right status line
set statusline+=\ %y	"Type of file in the buffer
set statusline+=\ (%l,%c)  "Line number and column number
set statusline+=\ %p%%\/%L	"Percent through the file and total line count
set statusline+=\ [%-3(%{SLFileSize()}%)]\	"File size of current buffer
set noshowmode	"Mode already handled in statusline

" ------------------------------ Leader Mappings -----------------------------
nmap <leader>p "+p|  "Paste from system clipboard
nmap <leader>l :set list!<CR>|	"Toggle listing whitespace chars
nmap <leader>s :s/\s\+$//e<CR>|  "Remove trailing white spaces in line
nmap <leader>S :%s/\s\+$//e<CR>|	"Remove all trailing white spaces in buffer
nmap <leader>f :Lexplore %:p:h<CR>|  "Open netrw file browser in current file directory
nmap <leader>F :Lexplore <CR>|	"Open netrw file browser in the home directory
nmap <leader>ff :Files<CR>|  "Fuzzy find files
nmap <leader>FF :Files $HOME<CR>|  "Fuzzy find files
nmap <leader>v :source $MYVIMRC<CR>|	"Source this .vimrc
nmap <leader>c :FormatLines<CR>|	"Format lines of code with google/vim-codefmt
nmap <leader>C :FormatCode<CR>|  "Format entire buffer with google/vim-codefmt

augroup run_code
	autocmd Filetype Bash nmap <leader>r :!bash %<CR>
	autocmd Filetype Python nmap <leader>r :!python3 %<CR>
augroup END

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" Autocommands
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
"Keep view options like folds extant automatically
augroup maintain_view
	autocmd BufWinLeave *.* mkview
	autocmd BufWinEnter *.* silent loadview 
augroup END

augroup filetypes  "Filetype-specific settings
	autocmd BufNewFile,BufRead *.md,*.txt  setlocal spell wrap linebreak
	autocmd Filetype python setlocal noexpandtab tabstop=2 shiftwidth=2 autoindent
	autocmd Filetype man setlocal nolist colorcolumn=""
augroup END

augroup commenting	"Filetype-specific comment strings
	autocmd FileType javascriptreact,typescriptreact setlocal commentstring={/*\ %s\ */}	"React
augroup END

augroup templates
	autocmd BufNewFile *.sh 0r $HOME/.vim/templates/skeleton.sh  "Bash skeleton
	autocmd BufNewFile *.py 0r $HOME/.vim/templates/skeleton.py  "Python skeleton
augroup END

""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" Plugins
""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" ----------------------------- Install vim-plug ------------------------------
if empty(glob('$HOME/.vim/autoload/plug.vim'))	"Conditional vim-plug installer
	silent !curl -fsSLo $HOME/.vim/autoload/plug.vim --create-dirs
		\ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
	autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
endif

" ---------------------------------- Plugins ----------------------------------
call plug#begin()
Plug 'tpope/vim-commentary'  "Commenting made easy
Plug 'tpope/vim-surround'  "Quoting/parenthesizing made easy
Plug 'jiangmiao/auto-pairs'  "Auto-close pairs
Plug 'tpope/vim-repeat'  "Repeat compatability with vim surround
Plug 'michaeljsmith/vim-indent-object'	"Indentation as a text object
Plug 'junegunn/fzf', { 'do': { -> fzf#install()  }	}  "Fuzzy finder install
Plug 'junegunn/fzf.vim'  "Fuzzy finder
Plug 'arzg/vim-colors-xcode'	"XCode colorscheme
Plug 'ap/vim-css-color'  "CSS Color highlighter
Plug 'google/vim-maktaba'  "Dependency of google/vim-codefmt
Plug 'google/vim-codefmt'  "Auto-format code
Plug 'google/vim-glaive'	"Dependency of google/vim-codefmt
Plug 'neoclide/coc.nvim', {'branch': 'release'}  "Code completion
Plug 'airblade/vim-gitgutter'  "Git status in the gutter
call plug#end()

" -------------------------- Plugin-Specific Settings -------------------------
colorscheme xcodedarkhc  "https://github.com/arzg/vim-colors-xcode (high contrast)

" ------------------------------- Code Completion -----------------------------
"Tab trigger completion
inoremap <silent><expr> <TAB>
	\ coc#pum#visible() ? coc#pum#next(1) :
	\ CheckBackspace() ? "\<Tab>" :
	\ coc#refresh()
inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"

"<CR> accept completion
inoremap <silent><expr> <CR> coc#pum#visible() ? coc#pum#confirm()
	\: "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"

function! CheckBackspace() abort
	let col = col('.') - 1
	return !col || getline('.')[col - 1]	=~# '\s'
endfunction
