package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

////////////////////////////
//         UTILS          //
////////////////////////////

// Download a file from this repo and return as byte array
func download(path) []byte {
	// Get request
	url = "https://raw.githubusercontent.com/williamwmarx/shell/main/" + path
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
	cmd := exec.Command(command)
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
	_, err := toml.Decode(tomlText, &packages)
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
	if commandExists("apt") {
		return packageManager{
			name: "apt",
			installCmd: "apt install -y"
			updateCmd: "apt update"
		}
	} else if commandExists("brew") {
		return packageManager{
			name: "brew",
			installCmd: "brew install"
			updateCmd: "brew upgrade"
		}
	} else if commandExists("dnf") {
		return packageManager{
			name: "dnf",
			installCmd: "dnf install -y"
			updateCmd: "dnf update"
		}
	} else if commandExists("pacman") {
		return packageManager{
			name: "pacman",
			installCmd: "pacman -S --no-confirm"
			updateCmd: "pacman -Syu"
		}
	}
}

var pm packageManager = getPackageManager()

// Handle the installing
func install(packageGroups ...string) {
	if len(packageGroups) > 0 {
		// Update package manager
		pm.update()
		
		// Gather packages to install
		var packagesToInstall []string
		for _, packageGroup := range packageGroups {
			if packageGroup == "core" {
				for _, p := range packages.Core {
					packagesToInstall = append(packagesToInstall, p.Name)
				}
			}
		}
	}
}
