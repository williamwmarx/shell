package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func tmuxConfig(temporary bool) []action {
	var actions []action

	// Check if tmux already exists
	tmuxAlreadyExists := commandExists("tmux")

	// Config base path
	tmuxConfBasePath := "~"
	if temporary {
		tmuxConfBasePath = "~/.shell.tmp"
	}

	// Install tmux if necessary
	if (temporary && !tmuxAlreadyExists) || !temporary {
		actions = append(actions, action{pm.updateCmd, "Updating package manager"})
		actions = append(actions, action{pm.installCommand(packages.packageByName("tmux")), "Installing tmux"})
	}

	// Get and save tmux.conf
	actions = append(actions, action{formatCurl("tmux/tmux.conf", tmuxConfBasePath+"/.tmux.conf"), "Saving tmux.conf"})

	// Create uninstall script if temporary install
	if temporary {
		uninstallCommands := ""
		if !tmuxAlreadyExists {
			uninstallCommands += pm.uninstallCommand(packages.packageByName("tmux")) + " && "
		}
		uninstallCommands += "rm -rf ~/.shell.tmp"
		actions = append(actions, action{uninstallCommands, "Saving uninstall script"})
	}

	return actions
}

// Install vim, vimrc, templates, and plugins
func vimConfig() []action {
	actions := []action{{pm.updateCmd, "Updating package manager"}}

	// Install Vim
	actions = append(actions, action{pm.installCommand(packages.packageByName("Vim")), "Installing vim"})

	// Get and save vimrc
	actions = append(actions, action{formatCurl("vim/vimrc", "~/.vimrc"), "Saving vimrc"})

	// Get and save template files
	dir, err := os.Open("vim/templates")
	if err != nil {
		log.Fatal(err)
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	// Write each file to ~/.vim/templates
	for _, file := range files {
		downloadURL := fmt.Sprintf("vim/templates/%s", file.Name())
		downloadPath := fmt.Sprintf("~/.vim/templates/%s", file.Name())
		curlSkeletonFile := formatCurl(downloadURL, downloadPath)
		actions = append(actions, action{fmt.Sprintf("mkdir -p ~/.vim/templates && %s", curlSkeletonFile), fmt.Sprintf("Saving %s", file.Name())})
	}

	// Install plugins
	actions = append(actions, action{"vim +PlugInstall +qall", "Installing vim plugins"})

	return actions
}

// Install vim for temporary use on another machine
func vanillaVimConfig(temporary bool) []action {
	var actions []action

	// Install Vim if necessary
	vimAlreadyExists := commandExists("vim")
	if !vimAlreadyExists {
		actions = append(actions, action{pm.updateCmd, "Updating package manager"})
		actions = append(actions, action{pm.installCommand(packages.packageByName("Vim")), "Installing vim"})
	}

	// Set vimrc path
	vimrcPath := "~/.vimrc"

	// Create .shell.tmp directory if temporary install
	if temporary {
		actions = append(actions, action{"mkdir -p ~/.shell.tmp", "Creating .shell.tmp directory"})
		vimrcPath = "~/.shell.tmp/vimrc"
	}

	// Get and save vimrc
	actions = append(actions, action{formatCurl("vim/vanillaVimrc", vimrcPath), "Saving vimrc"})

	// Create uninstall script if temporary install
	if temporary {
		uninstallCommands := ""
		if !vimAlreadyExists {
			uninstallCommands += pm.installCommand(packages.packageByName("Vim")) + " && "
		}
		uninstallCommands += "rm -rf ~/.shell.tmp"
		actions = append(actions, action{uninstallCommands, "Saving uninstall script"})
	}

	return actions
}

// Long command to install oh-my-zsh non-interactively
func ohmyzshInstallCmd() string {
	return "curl -fsSLo ~/install-oh-my-zsh.sh https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh && sh ~/install-oh-my-zsh.sh --unattended && rm ~/install-oh-my-zsh.sh"
}

// Install zsh, oh-my-zsh, zshrc and other zsh config files
func zshConfig() []action {
	actions := []action{{pm.updateCmd, "Updating package manager"}}
	actions = append(actions, action{pm.installCommand(packages.packageByName("Zsh")), "Installing zsh"})
	actions = append(actions, action{ohmyzshInstallCmd(), "Installing oh-my-zsh"})
	actions = append(actions, action{formatCurl("zsh/zshrc", "~/.zshrc"), "Saving zshrc"})
	actions = append(actions, action{formatCurl("zsh/aliases", "~/.aliases"), "Saving aliases"})
	actions = append(actions, action{formatCurl("zsh/functions", "~/.functions"), "Saving functions"})
	actions = append(actions, action{formatCurl("zsh/t3.zsh-theme", "~/.oh-my-zsh/themes/t3.zsh-theme"), "Saving t3.zsh-theme"})

	return actions
}

// Install zsh for temporary use on another machine
func vanillaZshConfig(temporary bool) []action {
	var actions []action

	// Install zsh if necessary
	zshAlreadyExists := commandExists("zsh")
	if !zshAlreadyExists {
		actions = append(actions, action{pm.updateCmd, "Updating package manager"})
		actions = append(actions, action{pm.installCommand(packages.packageByName("Zsh")), "Installing zsh"})
	}

	// Set zshrc path
	zshBasePath := "~/.shell"

	// Create .shell.tmp directory if temporary install
	if temporary {
		actions = append(actions, action{"mkdir -p ~/.shell.tmp", "Creating .shell.tmp directory"})
		zshBasePath = "~/.shell.tmp"
	} else {
		actions = append(actions, action{"mkdir -p ~/.shell", "Creating .shell directory"})
	}

	// Get and save .zshrc, .aliases, and .functions
	actions = append(actions, action{formatCurl("zsh/vanillaZshrc", zshBasePath+"/zshrc"), "Saving zshrc"})
	actions = append(actions, action{formatCurl("zsh/aliases", zshBasePath+"/aliases"), "Saving aliases"})
	actions = append(actions, action{formatCurl("zsh/functions", zshBasePath+"/functions"), "Saving functions"})

	// Create uninstall script if temporary install
	if temporary {
		uninstallCommands := ""
		if !zshAlreadyExists {
			uninstallCommands += pm.installCommand(packages.packageByName("Zsh")) + " && "
		}
		actions = append(actions, action{"rm -rf ~/.shell.tmp", "Saving uninstall script"})
	}

	return actions
}

// Full config/install
func fullConfig() []action {
	// First, update the package manaer
	actions := []action{{pm.updateCmd, "Updating package manager"}}

	// Install homebrew if necessary
	if runtime.GOOS == "darwin" && !commandExists("brew") {
		brewInstallCommand := "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
		actions = append(actions, action{brewInstallCommand, "Installing Homebrew"})
	}

	// Install packages
	actions = append(actions, installActions("Core", "Design")...)

	// If on Ubuntu, alias fd to fdfind
	if runtime.GOOS == "linux" {
		runCommand("ln -sf /usr/bin/fdfind /usr/local/bin/fd")
	}

	// Install GUI apps if on macOS
	if runtime.GOOS == "darwin" {
		actions = append(actions, installActions("GuiCore", "GuiDesign")...)
	}

	// Install oh-my-zsh
	actions = append(actions, action{ohmyzshInstallCmd(), "Installing oh-my-zsh"})

	// Clone this repo into home directory
	actions = append(actions, action{"git clone https://github.com/williamwmarx/shell.git ~/.shell", "Cloning shell repo"})

	// Create symlinks
	actions = append(actions, action{"ln -sf ~/.shell/git/gitconfig ~/.gitconfig", "Creating gitconfig symlink"})
	actions = append(actions, action{"mkdir -p ~/.gnupg && fd . ~/.shell/gnupg -x ln -sf {} ~/.gnupg/{/}", "Creating GPG symlinks"})
	actions = append(actions, action{"ln -sf ~/.shell/tmux/tmux.conf ~/.tmux.conf", "Creating tmux.conf symlink"})
	actions = append(actions, action{"ln -sf ~/.shell/vim/vimrc ~/.vimrc", "Creating vimrc symlink"})
	actions = append(actions, action{"mkdir -p ~/.vim/templates && fd . ~/.shell/vim/templates -x ln -sf {} ~/.vim/templates/{/}", "Creating vim skeleton symlinks"})
	actions = append(actions, action{"ln -sf ~/.shell/zsh/zshrc ~/.zshrc", "Creating zshrc symlink"})
	actions = append(actions, action{"ln -sf ~/.shell/zsh/aliases ~/.aliases", "Creating aliases symlink"})
	actions = append(actions, action{"ln -sf ~/.shell/zsh/functions ~/.functions", "Creating functions symlink"})
	actions = append(actions, action{"mkdir -p ~/.oh-my-zsh/themes && ln -sf ~/.shell/zsh/t3.zsh-theme ~/.oh-my-zsh/themes/t3.zsh-theme", "Creating t3.zsh-theme symlink"})

	// If on macOS, create macOS-specific symlinks
	if runtime.GOOS == "darwin" {
		actions = append(actions, action{"mkdir -p ~/.raycast && fd . ~/.shell/raycast -x ln -sf {} ~/.raycast{/}", "Creating Raycast shell script symlinks"})
		actions = append(actions, action{"ln -sf ~/.shell/skhd/skhdrc ~/.skhdrc", "Creating skhdrc symlink"})
		actions = append(actions, action{"ln -sf ~/.shell/yabai/yabairc ~/.yabairc", "Creating yabairc symlink"})
	}

	return actions
}
