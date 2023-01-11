package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"sort"
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

	// Ensure install dir exists
	runCommand("mkdir -p " + installDir)

	// Keep track of packages to uninstall if tmp install
	installFound := false
	var uninstallCommands []string

	// Iterate through install actions, formatting properly, and adding to actions
	var actions []action
	for _, a := range installer {
		// Get action parameters
		msg := a["msg"]
		cmd := a["cmd"]

		if strings.HasPrefix(cmd, "@install") {
			// Note that install was found
			installFound = true

			// Add package install action
			name := strings.TrimSpace(strings.TrimPrefix(cmd, "@install"))
			actions = append(actions, action{msg, PM.installCmd(name)})

			// Add uninstall command for package if --tmp passed
			if tmp {
				uninstallCommands = append(uninstallCommands, PM.uninstallCmd(name))
			}
		} else if strings.HasPrefix(cmd, "@save") {
			// Get file to save and turn into regex pattern — allows for wildcard matching
			filesToSave := strings.TrimSpace(strings.TrimPrefix(cmd, "@save"))
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
					var localPath string
					if tmp && contains(Config.Metadata.GitPaths, p) {
						// Get local parent dir of non-tmp install and add vanilla file name
						splitPath := strings.Split(p, "/")
						localPath = fmt.Sprintf("%s/%s%s", installDir, parentDir(p), splitPath[len(splitPath)-1])
					} else {
						if lp := Config.SyncTargets()[p]; lp != "" {
							// If file is in sync targets, use that path
							localPath = lp
						} else {
							// Otherwise, use repo path prepended with "~/.", assuming it's a dotfile in the root dir
							localPath = "~/." + p
						}
					}
					localPath = strings.ReplaceAll(localPath, "~", installDir)
					curl := fmt.Sprintf("curl -fsSLo %s %s", localPath, Config.Metadata.BaseURL+p)

					// If parent directory is not ~, ensure directory exists before running curl command
					if pd := parentDir(localPath); pd != installDir {
						curl = fmt.Sprintf("mkdir -p %s; %s", pd, curl)
					}

					matchedFiles = append(matchedFiles, curl)
				}
			}
			actions = append(actions, action{msg, strings.Join(matchedFiles, "; ")})
		}
	}

	// Create and add uninstall script if --tmp passed
	if len(uninstallCommands) > 0 {
		uninstallScript := strings.Join(uninstallCommands, "; ") + fmt.Sprintf("; rm -rf %s", installDir)
		uninstallAction := fmt.Sprintf("echo '%s' > %s/uninstall.sh", uninstallScript, installDir)
		actions = append(actions, action{"Adding uninstall script", uninstallAction})
	}

	// Prepend package manager update action if install was found
	if installFound {
		actions = append([]action{{"Updating package manager", PM.commands.updateCmd}}, actions...)
	}

	return actions
}

// Full config/install
func fullConfig() []action {
	// First, update the package manaer
	actions := []action{{"Updating package manager", PM.commands.updateCmd}}

	// Install homebrew if necessary
	if runtime.GOOS == "darwin" && !commandExists("brew") {
		brewInstallCommand := "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
		actions = append(actions, action{brewInstallCommand, "Installing Homebrew"})
	}

	// Install packages
	// Sort packages by group name, irrespective of case
	var packageGroups []string
	for group := range PM.Packages {
		packageGroups = append(packageGroups, group)
	}

	// Add package install actions and note requirements
	for _, packageGroup := range Sorted(packageGroups) {
		// Add packages actions
		for _, pia := range PM.packageInstallActions(packageGroup) {
			actions = append(actions, pia)
		}
	}

	// Clone this repo into home directory
	gitClone := fmt.Sprintf("git clone https://github.com/%s/%s.git ~/.%s", Config.Metadata.User, Config.Metadata.Repo, Config.Metadata.Repo)
	gitCloneMsg := fmt.Sprintf("Cloning github.com/%s/%s to ~/.%s", Config.Metadata.User, Config.Metadata.Repo, Config.Metadata.Repo)
	actions = append(actions, action{gitCloneMsg, gitClone})

	// Create symlinks for dotfiles
	var symlinkActions []action
	for repoPath, localPath := range Config.SyncTargets() {
		msg := fmt.Sprintf("Creating %s symlink", localPath)
		symlink := fmt.Sprintf("ln -sf ~/.%s/%s %s", Config.Metadata.Repo, repoPath, localPath)

		// Get parent directory of localPath
		splitLocalPath := strings.Split(localPath, "/")
		parentDir := strings.Join(splitLocalPath[:len(splitLocalPath)-1], "/")

		// If parent directory is not ~, ensure directory exists before creating symlink
		if parentDir != "~" {
			symlink = fmt.Sprintf("mkdir -p %s; %s", parentDir, symlink)
		}

		symlinkActions = append(symlinkActions, action{msg, symlink})
	}

	// Sort symlink actions by message, irrespective of case
	sort.SliceStable(symlinkActions, func(i, j int) bool {
		return strings.ToLower(symlinkActions[i].msg) < strings.ToLower(symlinkActions[j].msg)
	})

	// Add symlink actions to actions
	actions = append(actions, symlinkActions...)

	return actions
}
