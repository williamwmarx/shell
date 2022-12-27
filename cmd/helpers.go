package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

////////////////////////////
//         UTILS          //
////////////////////////////

var basePath string = "https://raw.githubusercontent.com/williamwmarx/shell/universalize/"

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
		Description  string
		URL 				 string `toml:"url"`
		AptName      string `toml:"apt_name"`
		BrewName     string `toml:"brew_name"`
		BrewCaskName string `toml:"brew_cask_name"`
		DnfName      string `toml:"dnf_name"`
		PacmanName   string `toml:"pacman_name"`
		InstallCommand string `toml:"install_command"`
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

var Packages pkgs = getPackages()

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
			for _, p := range Packages.packageGroup(pg) {
				actions = append(actions, action{pm.installCommand(p), fmt.Sprintf("Installing %s", p.Name)})
			}
		}
	}
	return actions
}
