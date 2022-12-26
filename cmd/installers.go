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

// Download a file from this repo and return as byte array
func download(path string) []byte {
	// Get request
	url := "https://raw.githubusercontent.com/williamwmarx/shell/main/" + path
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
func getCurlDownloadCommand(repoPath, outputPath string) string {
	return fmt.Sprintf("curl -fsSLo %s https://raw.githubusercontent.com/williamwmarx/shell/main/%s", outputPath, repoPath)
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

// Get a package by name
func getPackage(packageName string, packages pkgs) pkg {
	for _, packageGroup := range []map[string]pkg{packages.Core, packages.Design, packages.GuiCore, packages.GuiDesign} {
		for _, pack := range packageGroup {
			if pack.Name == packageName {
				return pack
			}
		}
	}
	return pkg{}
}

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

// Get packages from package group
func getPackagesFromGroup(packageGroup map[string]pkg) []pkg {
	var pacs []pkg
	for _, pack := range packageGroup {
		pacs = append(pacs, pack)
	}
	// Sort by alphabetical order, irrespective of case
	sort.Slice(pacs, func(i, j int) bool {
		return strings.ToLower(pacs[i].Name) < strings.ToLower(pacs[j].Name)
	})
	return pacs
}

// Get install command for a package
func getPackageInstallCmd(p pkg, pm packageManager) string {
	var packageName string
	switch pm.name {
	case "apt":
		packageName = p.AptName
	case "brew":
		if reflect.ValueOf(&p).Elem().FieldByName("BrewCaskName").String() != "" {
			packageName = "--cask " + p.BrewCaskName
		} else {
			packageName = p.BrewName
		}
	case "dnf":
		packageName = p.DnfName
	case "pacman":
		packageName = p.PacmanName
	}
	return fmt.Sprintf("%s %s", pm.installCmd, packageName)
}

// Get uninstall command for a package
func getPackageUninstallCmd(p pkg, pm packageManager) string {
	var packageName string
	switch pm.name {
	case "apt":
		packageName = p.AptName
	case "brew":
		if reflect.ValueOf(&p).Elem().FieldByName("BrewCaskName").String() != "" {
			packageName = "--cask " + p.BrewCaskName
		} else {
			packageName = p.BrewName
		}
	case "dnf":
		packageName = p.DnfName
	case "pacman":
		packageName = p.PacmanName
	}
	return fmt.Sprintf("%s %s", pm.uninstallCmd, packageName)
}

// Handle the installing
func installActions(packageGroups ...string) []action {
	packageInstallActions := []action{}

	if len(packageGroups) > 0 {
		// Gather packages to install
		var packagesToInstall []pkg
		for _, packageGroup := range packageGroups {
			switch packageGroup {
			case "Core":
				packagesToInstall = append(packagesToInstall, getPackagesFromGroup(packages.Core)...)
			case "Design":
				packagesToInstall = append(packagesToInstall, getPackagesFromGroup(packages.Design)...)
			case "GuiCore":
				packagesToInstall = append(packagesToInstall, getPackagesFromGroup(packages.GuiCore)...)
			case "GuiDesign":
				packagesToInstall = append(packagesToInstall, getPackagesFromGroup(packages.Design)...)
			default:
				log.Fatalf("Package group \"%s\" is not valid", packageGroup)
			}
		}

		// Configure package install actions
		for _, packToInstall := range packagesToInstall {
			packageInstallActions = append(packageInstallActions, action{getPackageInstallCmd(packToInstall, pm), fmt.Sprintf("Installing %s", packToInstall.Name)})
		}
	}

	return packageInstallActions
}

////////////////////////////
//	SPECIFIC INSTALLERS   //
////////////////////////////

func ohmyzshInstallCmd() string {
	return "curl -fsSLo ~/install-oh-my-zsh.sh https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh && sh ~/install-oh-my-zsh.sh --unattended && rm ~/install-oh-my-zsh.sh"
}

func tmuxConfig(temporary bool) []action {
	actionsToRun := []action{}
	// Check if tmux already exists
	tmuxAlreadyExists := commandExists("tmux")
	// Config base path
	tmuxConfBasePath := "~"
	if temporary {
		tmuxConfBasePath = "~/.shell.tmp"
	}
	// Install tmux if necessary
	if (temporary && !tmuxAlreadyExists) || !temporary {
		// Update package manager
		actionsToRun = append(actionsToRun, action{pm.updateCmd, "Updating package manager"})
		// Install tmux
		tmuxInstallCmd := getPackageInstallCmd(getPackage("tmux", packages), pm)
		actionsToRun = append(actionsToRun, action{tmuxInstallCmd, "Installing tmux"})
	}
	// Get and save tmux.conf
	curlTmuxConf := getCurlDownloadCommand("tmux/tmux.conf", tmuxConfBasePath+"/.tmux.conf")
	actionsToRun = append(actionsToRun, action{curlTmuxConf, "Saving tmux.conf"})
	// Create uninstall script if temporary install
	if temporary {
		uninstallCommands := ""
		if !tmuxAlreadyExists {
			uninstallCommands += getPackageUninstallCmd(getPackage("tmux", packages), pm) + " && "
		}
		uninstallCommands += "rm -rf ~/.shell.tmp"
		actionsToRun = append(actionsToRun, action{uninstallCommands, "Saving uninstall script"})
	}
	return actionsToRun
}

func vimConfig() []action {
	actionsToRun := []action{{pm.updateCmd, "Updating package manager"}}
	// Install Vim
	vimInstallCmd := getPackageInstallCmd(getPackage("Vim", packages), pm)
	actionsToRun = append(actionsToRun, action{vimInstallCmd, "Installing vim"})
	// Get and save vimrc
	curlVimrc := getCurlDownloadCommand("vim/vimrc", "~/.vimrc")
	actionsToRun = append(actionsToRun, action{curlVimrc, "Saving vimrc"})

	// Get and save template files
	dir, err := os.Open("vim/templates")
	if err != nil {
		log.Fatal(err)
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// Write each file to ~/.vim/templates
		curlSkeletonFile := getCurlDownloadCommand(fmt.Sprintf("vim/templates/%s", file.Name()), fmt.Sprintf("~/.vim/templates/%s", file.Name()))
		actionsToRun = append(actionsToRun, action{fmt.Sprintf("mkdir -p ~/.vim/templates && %s", curlSkeletonFile), fmt.Sprintf("Saving %s", file.Name())})
	}
	// Install plugins
	actionsToRun = append(actionsToRun, action{"vim +PlugInstall +qall", "Installing vim plugins"})
	return actionsToRun
}

func vanillaVimConfig(temporary bool) []action {
	actionsToRun := []action{}
	// Install Vim if necessary
	vimAlreadyExists := commandExists("vim")
	if !vimAlreadyExists {
		// Update package manager
		actionsToRun = append(actionsToRun, action{pm.updateCmd, "Updating package manager"})
		// Install Vim
		vimInstallCmd := getPackageInstallCmd(getPackage("Vim", packages), pm)
		actionsToRun = append(actionsToRun, action{vimInstallCmd, "Installing vim"})
	}
	// Set vimrc path
	vimrcPath := "~/.vimrc"
	// Create .shell.tmp directory if temporary install
	if temporary {
		actionsToRun = append(actionsToRun, action{"mkdir -p ~/.shell.tmp", "Creating .shell.tmp directory"})
		vimrcPath = "~/.shell.tmp/vimrc"
	}
	// Get and save vimrc
	curlVimrc := getCurlDownloadCommand("vim/vanillaVimrc", vimrcPath)
	actionsToRun = append(actionsToRun, action{curlVimrc, "Saving vimrc"})
	// Create uninstall script if temporary install
	if temporary {
		uninstallCommands := ""
		if !vimAlreadyExists {
			uninstallCommands += getPackageUninstallCmd(getPackage("Vim", packages), pm) + " && "
		}
		uninstallCommands += "rm -rf ~/.shell.tmp"
		actionsToRun = append(actionsToRun, action{uninstallCommands, "Saving uninstall script"})
	}
	return actionsToRun
}

func zshConfig() []action {
	actionsToRun := []action{{pm.updateCmd, "Updating package manager"}}
	// Install zsh
	zshInstallCmd := getPackageInstallCmd(getPackage("Zsh", packages), pm)
	actionsToRun = append(actionsToRun, action{zshInstallCmd, "Installing zsh"})
	// Install oh-my-zsh non-interactively
	actionsToRun = append(actionsToRun, action{ohmyzshInstallCmd(), "Installing oh-my-zsh"})
	// Get and save zshrc
	curlZshrc := getCurlDownloadCommand("zsh/zshrc", "~/.zshrc")
	actionsToRun = append(actionsToRun, action{curlZshrc, "Saving zshrc"})
	// Get and save aliases
	curlAliases := getCurlDownloadCommand("zsh/aliases", "~/.aliases")
	actionsToRun = append(actionsToRun, action{curlAliases, "Saving aliases"})
	// Get and save functions
	curlFunctions := getCurlDownloadCommand("zsh/functions", "~/.functions")
	actionsToRun = append(actionsToRun, action{curlFunctions, "Saving functions"})
	// Get and save t3.zsh-theme
	curlT3Theme := getCurlDownloadCommand("zsh/t3.zsh-theme", "~/.oh-my-zsh/themes/t3.zsh-theme")
	actionsToRun = append(actionsToRun, action{curlT3Theme, "Saving t3.zsh-theme"})
	return actionsToRun
}

func vanillaZshConfig(temporary bool) []action {
	actionsToRun := []action{}
	// Install zsh if necessary
	zshAlreadyExists := commandExists("zsh")
	if !zshAlreadyExists {
		// Update package manager
		actionsToRun = append(actionsToRun, action{pm.updateCmd, "Updating package manager"})
		// Install zsh
		zshInstallCmd := getPackageInstallCmd(getPackage("Zsh", packages), pm)
		actionsToRun = append(actionsToRun, action{zshInstallCmd, "Installing zsh"})
	}
	// Set zshrc path
	zshBasePath := "~/.shell"
	// Create .shell.tmp directory if temporary install
	if temporary {
		actionsToRun = append(actionsToRun, action{"mkdir -p ~/.shell.tmp", "Creating .shell.tmp directory"})
		zshBasePath = "~/.shell.tmp"
	} else {
		actionsToRun = append(actionsToRun, action{"mkdir -p ~/.shell", "Creating .shell directory"})
	}
	// Get and save zshrc
	curlZshrc := getCurlDownloadCommand("zsh/vanillaZshrc", zshBasePath+"/zshrc")
	actionsToRun = append(actionsToRun, action{curlZshrc, "Saving zshrc"})
	// Get and save aliases
	curlAliases := getCurlDownloadCommand("zsh/aliases", zshBasePath+"/aliases")
	actionsToRun = append(actionsToRun, action{curlAliases, "Saving aliases"})
	// Get and save functions
	curlFunctions := getCurlDownloadCommand("zsh/functions", zshBasePath+"/functions")
	actionsToRun = append(actionsToRun, action{curlFunctions, "Saving functions"})
	// Create uninstall script if temporary install
	if temporary {
		uninstallCommands := ""
		if !zshAlreadyExists {
			uninstallCommands += getPackageUninstallCmd(getPackage("Zsh", packages), pm) + " && "
		}
		uninstallCommands += "rm -rf ~/.shell.tmp"
		actionsToRun = append(actionsToRun, action{uninstallCommands, "Saving uninstall script"})
	}
	return actionsToRun
}

// Full config/install
func fullConfig() []action {
	// First, update the package manaer
	actionsToRun := []action{{pm.updateCmd, "Updating package manager"}}

	// Install homebrew if necessary
	if runtime.GOOS == "darwin" && !commandExists("brew") {
		brewInstallCommand := "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
		actionsToRun = append(actionsToRun, action{brewInstallCommand, "Installing Homebrew"})
	}

	// Install packages
	actionsToRun = append(actionsToRun, installActions("Core")...)
	actionsToRun = append(actionsToRun, installActions("Design")...)

	// If on Ubuntu, alias fd to fdfind
	if runtime.GOOS == "linux" {
		runCommand("ln -sf /usr/bin/fdfind /usr/local/bin/fd")
	}

	// Install GUI apps if on macOS
	if runtime.GOOS == "darwin" {
		actionsToRun = append(actionsToRun, installActions("GuiCore")...)
		actionsToRun = append(actionsToRun, installActions("GuiDesign")...)
	}

	// Install oh-my-zsh
	actionsToRun = append(actionsToRun, action{ohmyzshInstallCmd(), "Installing oh-my-zsh"})

	// Clone this repo into home directory
	actionsToRun = append(actionsToRun, action{"git clone https://github.com/williamwmarx/shell.git ~/.shell", "Cloning shell repo"})

	// Create symlinks
	symlinks := map[string]string{
		"ln -sf ~/.shell/git/gitconfig ~/.gitconfig":                                                        "Creating gitconfig symlink",
		"mkdir -p ~/.gnupg && fd . ~/.shell/gnupg -x ln -sf {} ~/.gnupg/{/}":                                "Creating GPG symlinks",
		"ln -sf ~/.shell/tmux/tmux.conf ~/.tmux.conf":                                                       "Creating tmux.conf symlink",
		"ln -sf ~/.shell/vim/vimrc ~/.vimrc":                                                                "Creating vimrc symlink",
		"mkdir -p ~/.vim/templates && fd . ~/.shell/vim/templates -x ln -sf {} ~/.vim/templates/{/}":        "Creating vim skeleton symlinks",
		"ln -sf ~/.shell/zsh/zshrc ~/.zshrc":                                                                "Creating zshrc symlink",
		"ln -sf ~/.shell/zsh/aliases ~/.aliases":                                                            "Creating aliases symlink",
		"ln -sf ~/.shell/zsh/functions ~/.functions":                                                        "Creating functions symlink",
		"mkdir -p ~/.oh-my-zsh/themes && ln -sf ~/.shell/zsh/t3.zsh-theme ~/.oh-my-zsh/themes/t3.zsh-theme": "Creating t3.zsh-theme symlink",
	}

	// If on macOS, create macOS-specific symlinks
	if runtime.GOOS == "darwin" {
		macOSSymlinks := map[string]string{
			"mkdir -p ~/.raycast && fd . ~/.shell/raycast -x ln -sf {} ~/.raycast{/}": "Creating Raycast shell script symlinks",
			"ln -sf ~/.shell/skhd/skhdrc ~/.skhdrc":                                   "Creating skhdrc symlink",
			"ln -sf ~/.shell/yabai/yabairc ~/.yabairc":                                "Creating yabairc symlink",
		}
		// Append macOSSymlinks to symlinks
		for symlink, message := range macOSSymlinks {
			symlinks[symlink] = message
		}
	}

	// Append symlinks to actionsToRun
	for symlink, message := range symlinks {
		actionsToRun = append(actionsToRun, action{symlink, message})
	}

	return actionsToRun
}
