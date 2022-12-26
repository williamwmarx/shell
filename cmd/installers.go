package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

////////////////////////////
//         UTILS          //
////////////////////////////

var basePath string = "https://raw.githubusercontent.com/williamwmarx/shell/main/"

// Download a file from this repo and return as byte array
func download(path string) []byte {
	// Get request
	url := basePath + path
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Return bytes
	return b
}

// Create curl download command text given relative path in this repo and output path
func formatCurl(repoPath, outputPath string) string {
	return fmt.Sprintf("curl -fsSLo %s %s%s", outputPath, basePath, repoPath)
}

// Check if a command exists on the system
func commandExists(commandName string) bool {
	_, err := exec.LookPath(commandName)
	return err == nil
}

// Run a system command
func runCommand(command string) {
	// Get home directory, and replace all instances of ~ with it
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	command = strings.ReplaceAll(command, "~", home)

	// Split commands by && and run each one
	for _, c := range strings.Split(command, "&&") {
		args := strings.Fields(strings.TrimSpace(c))
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

////////////////////////////
//      GET PACKAGES      //
////////////////////////////

// Struct to hold all package info
type (
	pkgs struct {
		Core      map[string]pkg
		Design    map[string]pkg
		GuiCore   map[string]pkg `toml:"gui_core"`
		GuiDesign map[string]pkg `toml:"gui_design"`
	}

	pkg struct {
		Name         string
		AptName      string `toml:"apt_name"`
		BrewName     string `toml:"brew_name"`
		BrewCaskName string `toml:"brew_cask_name"`
		DnfName      string `toml:"dnf_name"`
		PacmanName   string `toml:"pacman_name"`
	}
)

// Get a package by name
func (packs *pkgs) packageByName(name string) pkg {
	packageGroups := reflect.ValueOf(&packs).Elem().Elem()
	for i := 0; i < packageGroups.NumField(); i++ {
		p := packageGroups.Field(i).MapRange()
		for p.Next() {
			if p.Value().FieldByName("Name").String() == name {
				return p.Value().Interface().(pkg)
			}
		}
	}
	return pkg{}
}

// Get all packages in a group given its name
func (packs *pkgs) packageGroup(groupName string) []pkg {
	var packages []pkg
	// Get all packages in group
	p := reflect.ValueOf(&packs).Elem().Elem().FieldByName(groupName).MapRange()
	for p.Next() {
		packages = append(packages, p.Value().Interface().(pkg))
	}
	// Sort by alphabetical order, irrespective of case
	sort.Slice(packages, func(i, j int) bool {
		return strings.ToLower(packages[i].Name) < strings.ToLower(packages[j].Name)
	})
	return packages
}

// Get all packages from TOML file
func getPackages() pkgs {
	// Download TOML file from this repo
	tomlText := download("packages/packages.toml")

	// Unmarshal TOML file into struct
	var packages pkgs
	_, err := toml.Decode(string(tomlText), &packages)
	if err != nil {
		log.Fatal(err)
	}
	return packages
}

var packages pkgs = getPackages()

////////////////////////////
//    INSTALL PACKAGES    //
////////////////////////////

// Keep track of package manager commands
type packageManager struct {
	name         string
	installCmd   string
	uninstallCmd string
	updateCmd    string
}

// Get system package name
func (pm *packageManager) systemPackageName(pack pkg) string {
	var pName string
	switch pm.name {
	case "apt":
		pName = pack.AptName
	case "brew":
		if reflect.ValueOf(&pack).Elem().FieldByName("BrewCaskName").String() != "" {
			pName = "--cask " + pack.BrewCaskName
		} else {
			pName = pack.BrewName
		}
	case "dnf":
		pName = pack.DnfName
	case "pacman":
		pName = pack.PacmanName
	}
	return pName
}

// Get install command for a given package
func (pm *packageManager) installCommand(p pkg) string {
	return fmt.Sprintf("%s %s", pm.installCmd, pm.systemPackageName(p))
}

// Get uninstall command for a given package
func (pm *packageManager) uninstallCommand(p pkg) string {
	return fmt.Sprintf("%s %s", pm.uninstallCmd, pm.systemPackageName(p))
}

// Get system pacakge manager install command
func getPackageManager() packageManager {
	var pm packageManager
	if commandExists("pacman") {
		pm = packageManager{
			name:         "pacman",
			installCmd:   "pacman -S --no-confirm",
			uninstallCmd: "pacman -Rs --no-confirm",
			updateCmd:    "pacman -Syu",
		}
	} else if commandExists("dnf") {
		pm = packageManager{
			name:         "dnf",
			installCmd:   "dnf install -y",
			uninstallCmd: "dnf remove -y",
			updateCmd:    "dnf update",
		}
	} else if commandExists("brew") {
		pm = packageManager{
			name:         "brew",
			installCmd:   "brew install",
			uninstallCmd: "brew uninstall",
			updateCmd:    "brew upgrade",
		}
	} else if commandExists("apt") {
		pm = packageManager{
			name:         "apt",
			installCmd:   "apt install -y",
			uninstallCmd: "apt remove -y",
			updateCmd:    "apt update",
		}
	}
	return pm
}

var pm packageManager = getPackageManager()

// Get all install actions
func installActions(packageGroups ...string) []action {
	var actions []action
	if len(packageGroups) > 0 {
		for _, pg := range packageGroups {
			for _, p := range packages.packageGroup(pg) {
				actions = append(actions, action{pm.installCommand(p), fmt.Sprintf("Installing %s", p.Name)})
			}
		}
	}
	return actions
}

////////////////////////////
//	SPECIFIC INSTALLERS   //
////////////////////////////

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

func ohmyzshInstallCmd() string {
	return "curl -fsSLo ~/install-oh-my-zsh.sh https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh && sh ~/install-oh-my-zsh.sh --unattended && rm ~/install-oh-my-zsh.sh"
}

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
