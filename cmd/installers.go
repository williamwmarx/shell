package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// Return all actions for a given flag
func install(flag string, tmp bool) []action {
	// Get install directory (defaults to home), and replace all instances of ~ with it
	installDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Get install actions for flag from TOML config
	i := Config.Installers[flag]
	installer := i.Install
	if tmp {
		// Change install dir to if tmp install
		installDir = strings.ReplaceAll(Config.TmpDir, "@repo_name", Config.Metadata.Repo)

		// If there's a tmp install rule set, use that
		if i.TmpInstall != nil {
			installer = i.TmpInstall
		}
	}

	// Keep track of packages to uninstall if tmp install
	installFound := false
	var uninstallCommands []string

	// Iterate through install actions, formatting properly, and adding to actions
	var actions []action
	for msg, a := range installer {
		if strings.HasPrefix(a, "@install") {
			// Note that install was found
			installFound = true

			// Add package install action
			name := strings.TrimSpace(strings.TrimPrefix(a, "@install"))
			actions = append(actions, action{msg, PM.installCmd(name)})

			// Add uninstall command for package if --tmp passed
			if tmp {
				uninstallCommands = append(uninstallCommands, PM.uninstallCmd(name))
			}
		} else if strings.HasPrefix(a, "@save") {
			// Get file to save and turn into regex pattern — allows for wildcard matching
			filesToSave := strings.TrimSpace(strings.TrimPrefix(a, "@save"))
			regexpPatttern := strings.ReplaceAll(filesToSave, "*", ".*")

			// Find matches
			var matchedFiles []string
			for _, p := range Config.Metadata.GitPaths {
				// Check if file matches pattern
				matched, err := regexp.Match(regexpPatttern, []byte(p))
				if err != nil {
					log.Fatal(err)
				}
				if matched {
					// If file matches, add properly formatted curl command to download it
					localPath := Config.SyncTargets()[strings.ReplaceAll(p, "vanilla_", "")]
					localPath = strings.ReplaceAll(localPath, "~", installDir)
					curl := fmt.Sprintf("curl -fsSLo %s %s", localPath, Config.Metadata.BaseURL+p)
					matchedFiles = append(matchedFiles, curl)
				}
			}
			actions = append(actions, action{msg, strings.Join(matchedFiles, "; ")})
		}
	}

	// Create and add uninstall script if --tmp passed
	if len(uninstallCommands) > 0 {
		uninstallScript := strings.Join(uninstallCommands, "\n") + fmt.Sprintf("rm -rf %s", installDir)
		uninstallAction := fmt.Sprintf("echo '%s' > %s/uninstall.sh", uninstallScript, installDir)
		actions = append(actions, action{"Adding uninstall script", uninstallAction})
	}

	// Prepend package manager update action if install was found
	if installFound {
		actions = append([]action{{PM.commands.updateCmd, "Updating package manager"}}, actions...)
	}

	return actions
}

// Full config/install
func fullConfig() []action {
	// First, update the package manaer
	actions := []action{{PM.commands.updateCmd, "Updating package manager"}}

	// Install homebrew if necessary
	if runtime.GOOS == "darwin" && !commandExists("brew") {
		brewInstallCommand := "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
		actions = append(actions, action{brewInstallCommand, "Installing Homebrew"})
	}

	// Install packages
	for _, pkgGroup := range PM.packages {
		for packageName := range pkgGroup {
			// Ignore description, as it's not a package
			if packageName != "description" {
				// Get install command for package and add to actions if it exists
				installCommmand := PM.installCmd(packageName)
				if installCommmand != "" {
					actions = append(actions, action{installCommmand, "Installing " + packageName})
				}
			}
		}
	}

	// Clone this repo into home directory
	gitClone := fmt.Sprintf("git clone https://github.com/%s/%s.git ~/.%s", Config.Metadata.User, Config.Metadata.Repo, Config.Metadata.Repo)
	gitCloneMsg := fmt.Sprintf("Cloning %s/%s repo to ~/.%s", Config.Metadata.User, Config.Metadata.Repo, Config.Metadata.Repo)
	actions = append(actions, action{gitClone, gitCloneMsg})

	// Create symlinks for dotfiles
	for repoPath, localPath := range Config.SyncTargets() {
		msg := fmt.Sprintf("Creating %s symlink", localPath)
		actions = append(actions, action{fmt.Sprintf("ln -sf %s %s", repoPath, localPath), msg})
	}

	return actions
}
