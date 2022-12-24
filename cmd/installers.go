package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"runtime"
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

	// Read text from body and return
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

// Check if a command exists on the system
func commandExists(commandName string) bool {
	_, err := exec.LookPath(commandName)
	return err == nil
}

// Run a system command
func runCommand(command string) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
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
	name       string
	installCmd string
	updateCmd  string
}

// Install a package
func (pm *packageManager) install(packageName string) {
	runCommand(fmt.Sprintf("%s %s", pm.installCmd, packageName))
}

// Update pacakge manager
func (pm *packageManager) update() {
	runCommand(pm.updateCmd)
}

// Get system pacakge manager install command
func getPackageManager() packageManager {
	var pm packageManager
	if commandExists("pacman") {
		pm = packageManager{
			name:       "pacman",
			installCmd: "pacman -S --no-confirm",
			updateCmd:  "pacman -Syu",
		}
	} else if commandExists("dnf") {
		pm = packageManager{
			name:       "dnf",
			installCmd: "dnf install -y",
			updateCmd:  "dnf update",
		}
	} else if commandExists("brew") {
		pm = packageManager{
			name:       "brew",
			installCmd: "brew install",
			updateCmd:  "brew upgrade",
		}
	} else if commandExists("apt") {
		pm = packageManager{
			name:       "apt",
			installCmd: "apt install -y",
			updateCmd:  "apt update",
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

// Handle the installing
func installActions(packageGroups ...string) []action {
	packageInstallActions := []action{}

	if len(packageGroups) > 0 {
		// Update package manager
		pm.update()

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
				log.Fatal(fmt.Sprintf("Package group \"%s\" is not valid", packageGroup))
			}
		}

		// Configure package install actions
		for _, packToInstall := range packagesToInstall {
			systemPackageName := getPackageInstallCmd(packToInstall, pm)
			packageInstallActions = append(packageInstallActions, action{fmt.Sprintf("%s %s", pm.installCmd, systemPackageName), fmt.Sprintf("Installing %s", packToInstall.Name)})
		}
	}

	return packageInstallActions
}

// Full config/install
func fullConfig() []action {
	actionsToRun := []action{}
	// Install homebrew if necessary
	if runtime.GOOS == "darwin" && !commandExists("brew") {
		brewInstallCommand := "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
		actionsToRun = append(actionsToRun, action{brewInstallCommand, "Installing Homebrew"})
	}

	// Install packages
	actionsToRun = append(actionsToRun, installActions("Core")...)
	actionsToRun = append(actionsToRun, installActions("Design")...)
	actionsToRun = append(actionsToRun, installActions("GuiCore")...)
	actionsToRun = append(actionsToRun, installActions("GuiDesign")...)

	// Install oh-my-zsh
	actionsToRun = append(actionsToRun, action{"sh -c \"$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\" \"\" --unattended", "Installing oh-my-zsh"})

	// Clone this repo into home directory
	actionsToRun = append(actionsToRun, action{"git clone https://github.com/williamwmarx/shell ~/.shell", "Cloning shell repo"})

	// Create symlinks
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/git/gitconfig ~/.gitconfig", "Creating gitconfig symlink"})
	actionsToRun = append(actionsToRun, action{"mkdir -p ~/.gnupg && ln -s ~/.shell/gnupg/gpg-agent.conf ~/.gnupg/gpg-agent.conf", "Creating gpg-agent.conf symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/gnupg/gpg.conf ~/.gnupg/gpg.conf", "Creating gpg.conf symlink"})
	actionsToRun = append(actionsToRun, action{"mkdir -p ~/.raycast && ln -s ~/.raycast/clear-format ~/.raycast/clear-format", "Creating Raycast clear-format script symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.raycast/expand-url ~/.raycast/expand-url", "Creating Raycast expand-url script symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/skhd/skhdrc ~/.skhdrc", "Creating skhdrc symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/tmux/tmux.conf ~/.tmux.conf", "Creating tmux.conf symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/vim/vimrc ~/.vimrc", "Creating vimrc symlink"})
	actionsToRun = append(actionsToRun, action{"mkdir -p ~/.vim/templates && ln -s ~/.shell/vim/tempaltes/skeleton.py ~/.vim/templates/skeleton.py && ln -s ~/.shell/vim/tempaltes/skeleton.sh ~/.vim/templates/skeleton.sh", "Creating vim skeleton symlinks"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/yabai/yabairc ~/.yabairc", "Creating yabairc symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/zsh/zshrc ~/.zshrc", "Creating zshrc symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/zsh/aliases ~/.aliases", "Creating aliases symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/zsh/functions ~/.functions", "Creating functions symlink"})
	actionsToRun = append(actionsToRun, action{"ln -s ~/.shell/zsh/t3.zsh-theme ~/.oh-my-zsh/themes/t3.zsh-theme", "Creating t3.zsh-theme symlink"})

	return actionsToRun
}

////////////////////////////
//	SPECIFIC INSTALLERS   //
////////////////////////////

func tmuxConfig() []action {
	actionsToRun := []action{}
	// Install tmux
	tmuxInstallCmd := getPackageInstallCmd(getPackage("tmux", packages), pm)
	actionsToRun = append(actionsToRun, action{tmuxInstallCmd, "Installing tmux"})
	// Get and save tmux.conf
	tmuxConf := download("tmux/tmux.conf")
	actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.tmux.conf", tmuxConf), "Saving tmux.conf"})
	return actionsToRun
}

func vimConfig() []action {
	actionsToRun := []action{}
	// Install Vim
	vimInstallCmd := getPackageInstallCmd(getPackage("vim", packages), pm)
	actionsToRun = append(actionsToRun, action{vimInstallCmd, "Installing vim"})
	// Get and save vimrc
	vimrc := download("vim/vimrc")
	actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.vimrc && mkdir -p ~/.vim", vimrc), "Saving vimrc"})
	// Get and save template files
	files, _ := ioutil.ReadDir("vim/templates")
	for _, file := range files {
		// Write each file to ~/.vim/templates
		skeletonFile := download(fmt.Sprintf("vim/templates/%s", file.Name()))
		actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.vim/templates/%s", skeletonFile, file.Name()), fmt.Sprintf("Saving %s", file.Name())})
	}
	// Install plugins
	actionsToRun = append(actionsToRun, action{"vim +PlugInstall +qall", "Installing vim plugins"})
	return actionsToRun
}

func zshConfig() []action {
	actionsToRun := []action{}
	// Install zsh
	zshInstallCmd := getPackageInstallCmd(getPackage("zsh", packages), pm)
	actionsToRun = append(actionsToRun, action{zshInstallCmd, "Installing zsh"})
	// Install oh-my-zsh non-interactively
	actionsToRun = append(actionsToRun, action{"sh -c \"$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\" \"\" --unattended", "Installing oh-my-zsh"})
	// Get and save zshrc
	zshrc := download("zsh/zshrc")
	actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.zshrc", zshrc), "Saving zshrc"})
	// Get and save aliases
	aliases := download("zsh/aliases")
	actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.aliases", aliases), "Saving aliases"})
	// Get and save functions
	functions := download("zsh/functions")
	actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.functions", functions), "Saving functions"})
	// Get and save t3.zsh-theme
	t3Theme := download("zsh/t3.zsh-theme")
	actionsToRun = append(actionsToRun, action{fmt.Sprintf("echo \"%s\" > ~/.oh-my-zsh/themes/t3.zsh-theme", t3Theme), "Saving t3.zsh-theme"})
	return actionsToRun
}
