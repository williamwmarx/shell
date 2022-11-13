#!/bin/sh
################################################################################
# Startup
################################################################################
# Ask for the admin password upfront and keep alive until `install.sh` has finished
sudo -v
while true; do sudo -n true; sleep 60; kill -0 "$$" || exit; done 2>/dev/null &


################################################################################
# Helper Functions
################################################################################
# ------------------------------ Print Statements -----------------------------
# Different styles so I don't have to type ANSI escape codes each time
print_red() { echo "\033[31m$1\033[0m"; }
print_green() { echo "\033[32m$1\033[0m"; }
print_yellow() { echo "\033[33m$1\033[0m"; }
print_bold() { echo "\033[1m$1\033[0m"; }

# Consistent installing, successful install, and failed install messages
print_i() { print_yellow "Installing $1..."; }
print_i_s() { print_green "$1 successfully installated."; }
print_i_f() { print_red "$1 could not be installated."; }

# Successful input message if last jobs status is 0, else fail message
print_status() {
	LAST_JOB_STATUS=$1
	APP_NAME=$2
	if [[ $LAST_JOB_STATUS -eq 0 ]]; then
		print_i_s $APP_NAME;
	else
		print_i_f $APP_NAME;
	fi
}

# ------------------------------- Install Helpers -----------------------------
# Homebrew requires Xcode command line tools
install_xcode_tools() {
	print_i "Xcode Command Line Tools"
	xcode-select --install
	print_status $? "Xcode Command Line Tools"
}

# Homebrew installer (https://brew.sh)
install_homebrew() {
	NONINTERACTIVE=1 \
		/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
}

# Tools I use non-stop and want on my standard account
install_daily_tools() {
	if [ `uname` = "Darwin" ]; then
		print_i "daily tools"
		brew install bat curl exa fd fzf git gnupg m-cli openssl pinentry-mac ripgrep tmux vim zsh
		sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
		brew install --cask cleanmymac dropbox iina little-snitch mullvadvpn parallels raycast zoom
		print_status $? "daily tools"
	fi
}

# Tools I use for design and want installed under a separate user
# The crux of this second account's utility is in remove the annoyance of Adobe software running
# in my main user space. Hopefullly can dockerize this in the future.
install_design_tools() {
	if [ `uname` = "Darwin" ]; then
		print_i "design tools"
		brew_isntall ffmpeg
		brew_install --cask adobe-creative-cloud glyphs rhino cycling74-max
		print_status $? "design tools"
	fi
}


################################################################################
# Installation
################################################################################
# ------------------------------- macOS Install -------------------------------
if [ `uname` = "Darwin" ]; then
	# If we don't have homebrew installed, install it
	if [ ! `command -v brew` ]; then
		[ ! `command -v xcode-select` ] && install_xcode_tools
		install_homebrew
		brew analytics off  # Opt out of anonymous analytics brew uses
	else
		print_green "Homebrew already installed"
	fi

	# Choose install option do the installing
	print_bold "Select install pattern:"
	printf "  1: Daily tools\n  2: Design tools\n  3: All tools\nSelection: "
	read INSTALL_CHOICE
	case $INSTALL_CHOICE in
		1) install_daily_tools;;
		2) install_design_tools;;
		3) install_daily_tools; install_design_tools;;
		*) print_red "Invalid condition selected"
	esac

	# Clone this repo to ~/.dotfiles
	git clone https://github.com/williamwmarx/.dotfiles.git $HOME/.dotfiles

	# Create smylinks to dotfiles in this repo
	ln -s $HOME/.dotfiles/zsh/.zshrc $HOME/.zshrc
	ln -s $HOME/.dotfiles/zsh/.aliases $HOME/.aliases
	ln -s $HOME/.dotfiles/zsh/.functions $HOME/.functions
	ln -s $HOME/.dotfiles/zsh/t3.zsh-theme $HOME/.oh-my-zsh/themes/t3.zsh-theme
	ln -s $HOME/.dotfiles/vim/.vimrc $HOME/.vimrc
	ln -s $HOME/.dotfiles/vim/.vim $HOME/.vim

	# Install Vim plugins
	vim -c "PlugInstall" -c "q" -c "q"

	# Open Terminal theme
	# Need a more elegant solution here
	open $HOME/.dotfiles/macOS/t3.terminal

	# Update macOS defaults
	sudo bash $HOME/.dotfiles/macOS/.macos

	# Reboot to allow all changes to take place
	sudo shutdown -r now
fi
