package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
)

///////////////////////////
//     GENERAL UTILS     //
///////////////////////////

// Check if a command exists on the system
func commandExists(commandName string) bool {
	_, err := exec.LookPath(commandName)
	return err == nil
}

// Run a system command and get output
func runCommand(command string) string {
	cmd := exec.Command("sh", "-c", command)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(stdout)
}

// Download a file and return as byte array
func download(url string) []byte {
	// Get request
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

// Check if string is contained in array of strings
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Get parent directory of path
func parentDir(path string) string {
	splitLocalPath := strings.Split(path, "/")
	return strings.Join(splitLocalPath[:len(splitLocalPath)-1], "/")
}

// Sort an array of strings, irrespective of case
func Sorted(s []string) []string {
	sort.Slice(s, func(i, j int) bool {
		return strings.ToLower(s[i]) < strings.ToLower(s[j])
	})
	return s
}

///////////////////////////
//        CONFIG         //
///////////////////////////

// Structs to store contents of config.toml
type (
	config struct {
		TmpDir          string `toml:"tmp_dir"`
		InstallURL      string `toml:"custom_install_url"`
		HelpDescription string `toml:"help_description"`
		Sync            map[string]targetClass
		Installers      map[string]Installer
		Metadata        metadata
	}

	targetClass struct {
		Name      string
		MacOSOnly bool `toml:"macos_only"`
		Targets   []Target
	}

	Target struct {
		Description string
		RepoPath    string `toml:"repo_path"`
		LocalPath   string `toml:"local_path"`
	}

	Installer struct {
		HelpMessage string `toml:"help_message"`
		Description string
		Install     []map[string]string
		TmpInstall  []map[string]string `toml:"tmp_install"`
	}

	metadata struct {
		User     string
		Repo     string
		BaseURL  string
		GitPaths []string
	}
)

// Create map of repo paths and local paths for all sync targets
func (c *config) SyncTargets() map[string]string {
	targets := make(map[string]string)
	for _, s := range Config.Sync {
		for _, t := range s.Targets {
			targets[t.RepoPath] = t.LocalPath
		}
	}
	return targets
}

// Get all paths in remote GitHub repo
func remoteGitPaths(user, repo, branch string) []string {
	// Get tree from GitHub API
	client := github.NewClient(nil)
	tree, _, err := client.Git.GetTree(context.Background(), user, repo, branch, true)
	if err != nil {
		log.Fatal(err)
	}

	// Get paths from tree
	var gitPaths []string
	for _, entry := range tree.Entries {
		gitPaths = append(gitPaths, entry.GetPath())
	}

	return gitPaths
}

// Unmarshall config.toml file and add metadata
func getConfig() config {
	var c config

	// Get git repo remote origin url and split into user and repo
	// We need to do this first to get the remote config file
	gitRemoteOriginURL := "https://github.com/williamwmarx/shell"
	splitURL := strings.Split(strings.TrimSpace(gitRemoteOriginURL), "/")
	c.Metadata.User = splitURL[len(splitURL)-2]
	c.Metadata.Repo = splitURL[len(splitURL)-1]

	// Get base URL for raw GitHub user content
	branch := "main"
	c.Metadata.BaseURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/", c.Metadata.User, c.Metadata.Repo, branch)

	// Read text of TOML file
	configToml := download(c.Metadata.BaseURL + "config.toml")

	// Unmarshall TOML file
	_, err := toml.Decode(string(configToml), &c)
	if err != nil {
		log.Fatal(err)
	}

	// Update TmpDir with repo info
	c.TmpDir = strings.ReplaceAll(c.TmpDir, "@repo_name", c.Metadata.Repo)

	// If no custom install url, set install url
	if c.InstallURL == "" {
		c.InstallURL = c.Metadata.BaseURL + "install.sh"
	}

	// Get all paths in remote GitHub repo
	c.Metadata.GitPaths = remoteGitPaths(c.Metadata.User, c.Metadata.Repo, branch)

	return c
}

// Make config data globally available
var Config config = getConfig()

///////////////////////////
//    PACKAGE MANAGER    //
///////////////////////////

// Store package manager commands and available packages
type (
	packageManager struct {
		commands pmCommands
		Packages pkgGroup
	}

	pmCommands struct {
		name         string
		installCmd   string
		uninstallCmd string
		updateCmd    string
	}

	pkgGroup map[string]pkgs

	pkgs struct {
		Description string
		Packages    map[string]map[string]string
	}
)

// Get a package by its name
func (p *pkgGroup) PackageByName(name string) map[string]string {
	for _, group := range *p {
		if pack, packInGroup := group.Packages[name]; packInGroup {
			return pack
		}
	}
	return map[string]string{}
}

// Get system install command for a given package
func (pm *packageManager) installCmd(name string) string {
	// Get package from packages.toml
	pack := pm.Packages.PackageByName(name)

	// Check if pack has key "install_command" and return it if it does
	if installCmd, ok := pack["install_command"]; ok {
		return installCmd
	}

	// Check for system package name and return install command if it exists
	if systemPackageName, ok := pack[pm.commands.name]; ok {
		return pm.commands.installCmd + " " + systemPackageName
	}

	// Package not found, return empty string
	return ""
}

// Get system uninstall command for a given package
func (pm *packageManager) uninstallCmd(name string) string {
	// Get package from packages.toml
	pack := pm.Packages.PackageByName(name)

	// Check if pack has key "uninstall_command" and return it if it does
	if installCmd, ok := pack["uninstall_command"]; ok {
		return installCmd
	}

	// Check for system package name and return uninstall command if it exists
	if systemPackageName, ok := pack[pm.commands.name]; ok {
		return pm.commands.uninstallCmd + " " + systemPackageName
	}

	// Package not found, return empty string
	return ""
}

// Type for package install actions
type packageAction struct {
	a        action
	requires string
}

// Get system install commands for a given package group
func (pm *packageManager) packageInstallActions(packageGroupName string) []action {
	// Sort packageNames by name, irrespective of case
	var packageNames []string
	for packageName := range pm.Packages[packageGroupName].Packages {
		packageNames = append(packageNames, packageName)
	}

	// Add package install commands
	var packageActions []packageAction
	for _, packageName := range Sorted(packageNames) {
		// Ignore description, as it's not a package
		if packageName != "description" {
			// Get install command for package and add to actions if it exists
			installCommand := pm.installCmd(packageName)
			if installCommand != "" {
				// Get requirement for package
				var requires string
				if r, ok := pm.Packages.PackageByName(packageName)["requires"]; ok {
					requires = r
				}

				// Add package install action to packageActions
				a := action{"Installing " + packageName, installCommand}

				packageActions = append(packageActions, packageAction{a, requires})
			}
		}
	}

	// Add package actions to actions, ensuring packages that require others are installed after their dependencies
	var actions []action
	var packagesWithDependencies []action
	for _, pa := range packageActions {
		if pa.requires != "" {
			packagesWithDependencies = append(packagesWithDependencies, pa.a)
		} else {
			actions = append(actions, pa.a)
		}
	}

	// Add packages with dependencies to actions
	actions = append(actions, packagesWithDependencies...)

	return actions
}

// Get system pacakge manager commands and listed packages
func getPackageManager() packageManager {
	var pm packageManager

	// Get package manager commands
	if commandExists("pacman") {
		pm.commands = pmCommands{
			name:         "pacman",
			installCmd:   "pacman -S --no-confirm",
			uninstallCmd: "pacman -Rs --no-confirm",
			updateCmd:    "pacman -Syu",
		}
	} else if commandExists("dnf") {
		pm.commands = pmCommands{
			name:         "dnf",
			installCmd:   "dnf install -y",
			uninstallCmd: "dnf remove -y",
			updateCmd:    "dnf update",
		}
	} else if commandExists("brew") {
		pm.commands = pmCommands{
			name:         "brew",
			installCmd:   "brew install",
			uninstallCmd: "brew uninstall",
			updateCmd:    "brew upgrade",
		}
	} else if commandExists("apt") {
		pm.commands = pmCommands{
			name:         "apt",
			installCmd:   "apt install -y",
			uninstallCmd: "apt remove -y",
			updateCmd:    "apt update",
		}
	}

	// Download packages TOML file from this repo
	tomlText := download(Config.Metadata.BaseURL + "packages.toml")

	// Unmarshal TOML file into struct
	_, err := toml.Decode(string(tomlText), &pm.Packages)
	if err != nil {
		log.Fatal(err)
	}

	return pm
}

// Make package manager and packages available globally
var PM packageManager = getPackageManager()
