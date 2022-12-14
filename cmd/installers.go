package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
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

// Handle the installing
func install(packageGroups ...string) {
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

		// Do the installing
		for _, packToInstall := range packagesToInstall {
			var packageName string
			switch pm.name {
			case "apt":
				packageName = packToInstall.AptName
			case "brew":
				if reflect.ValueOf(&packToInstall).Elem().FieldByName("BrewCaskName").String() != "" {
					packageName = "--cask " + packToInstall.BrewCaskName
				} else {
					packageName = packToInstall.BrewName
				}
			case "dnf":
				packageName = packToInstall.DnfName
			case "pacman":
				packageName = packToInstall.PacmanName
			}

			runCommand(fmt.Sprintf("%s %s", pm.installCmd, packageName))
		}
	}
}

// Full config/install
func fullConfig(b bool) {
	// Install homebrew if necessary
	if (runtime.GOOS == "darwin" && !commandExists("brew")) {
		runCommand("NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	}

	// Install packages
	install("Core")
	install("Design")
	install("GuiCore")
	install("GuiDesign")

	// Install oh-my-zsh
	runCommand("sh -c \"$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\" \"\" --unattended")

	// Clone this repo into home directory
	runCommand("git clone https://github.com/williamwmarx/shell ~/.shell")

	// Create symlinks
	os.Symlink("~/.shell/git/gitconfig", "~/.gitconfig")
	_ = os.Mkdir("~/.gnupg", os.ModePerm)
	os.Symlink("~/.shell/gnupg/gpg-agent.conf", "~/.gnupg/gpg-agent.conf")
	os.Symlink("~/.shell/gnupg/gpg.conf", "~/.gnupg/gpg.conf")
	os.Symlink("~/.shell/personal/plan", "~/.plan")
	_ = os.Mkdir("~/.raycast", os.ModePerm)
	os.Symlink("~/.shell/raycast/clear-format.sh", "~/.raycast/clear-format.sh")
	os.Symlink("~/.shell/raycast/expand-url.sh", "~/.raycast/expand-url.sh")
	os.Symlink("~/.shell/skhd/skhdrc", "~/.skhdrc")
	os.Symlink("~/.shell/tmux/tmux.conf", "~/.tmux.conf")
	os.Symlink("~/.shell/vim/vimrc", "~/.vimrc")
	_ = os.Mkdir("~/.vim", os.ModePerm)
	_ = os.Mkdir("~/.vim/templates", os.ModePerm)
	os.Symlink("~/.shell/vim/tempaltes/skeleton.py", "~/.vim/templates/skeleton.py")
	os.Symlink("~/.shell/vim/tempaltes/skeleton.sh", "~/.vim/templates/skeleton.sh")
	os.Symlink("~/.shell/yabai/yabairc", "~/.yabairc")
	os.Symlink("~/.shell/zsh/zshrc", "~/.zshrc")
	os.Symlink("~/.shell/zsh/aliases", "~/.aliases")
	os.Symlink("~/.shell/zsh/functions", "~/.functions")
	os.Symlink("~/.shell/zsh/t3.zsh-theme", "~/.oh-my-zsh/themes/t3.zsh-theme")
}
