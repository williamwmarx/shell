package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/williamwmarx/shell/cmd"
)

// Read file to string
func readToString(path string) string {
	markdown, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(markdown)
}

// Write string to file
func writeString(path, markdown string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(markdown)
}

// Writes apex README.md, showing install command and listing synced dotfiles
func writeREADME() {
	// Read the apex README template and insert proper URL
	markdown := readToString("assembly/README_TEMPLATE.md")
	markdown = strings.ReplaceAll(markdown, "%INSTALL_URL%", cmd.Config.InstallURL)

	// Format text for synced targets
	var syncTargets []string
	for _, s := range cmd.Config.Sync {
		// Format header
		header := fmt.Sprintf("### %s\n\n", s.Name)
		// Format targets
		var targets []string
		for _, t := range s.Targets {
			path := strings.Split(t.LocalPath, "/")
			target := fmt.Sprintf("- [%s](%s) — %s\n", path[len(path)-1], t.RepoPath, t.Description)
			targets = append(targets, target)
		}
		// Add header and sorted targets to syncTargetsText
		syncTargets = append(syncTargets, header+strings.Join(cmd.Sorted(targets), "")+"\n")
	}

	// Add dotfiles to markdown
	dotfilesText := strings.TrimSpace(strings.Join(cmd.Sorted(syncTargets), ""))
	markdown = strings.ReplaceAll(markdown, "%DOTFILES%", dotfilesText)

	var packages []string
	for packageGroup := range cmd.PM.Packages {
		header := fmt.Sprintf("### %s\n\n", packageGroup)
		header += cmd.PM.Packages[packageGroup].Description + "\n\n"
		var pgPackages []string
		for pName, p := range cmd.PM.Packages[packageGroup].Packages {
			pgPackages = append(pgPackages, fmt.Sprintf("- [%s](%s) - %s\n", pName, p["url"], p["description"]))
		}
		packages = append(packages, header+strings.Join(cmd.Sorted(pgPackages), "")+"\n")
	}

	// Add packages to markdown
	packagesText := strings.TrimSpace(strings.Join(cmd.Sorted(packages), ""))
	markdown = strings.ReplaceAll(markdown, "%PACKAGES%", packagesText)

	// Write markdown to README.md
	writeString("README.md", markdown)
}

// Writes INSTALL.md, showing thorough install options
func writeINSTALL() {
	// Read the apex README template and insert proper URL
	markdown := readToString("assembly/INSTALL_TEMPLATE.md")
	markdown = strings.ReplaceAll(markdown, "%INSTALL_URL%", cmd.Config.InstallURL)

	// Get installers and sort by name
	var installerNames []string
	for k := range cmd.Config.Installers {
		installerNames = append(installerNames, k)
	}

	// Add installers to markdown
	var installers []string
	for _, inst := range installerNames {
		// Title and description
		instText := fmt.Sprintf("#### %s\n\n%s\n\n", inst, cmd.Config.Installers[inst].Description)
		// Code block
		instText += fmt.Sprintf("```bash\nsh <(curl %s) --%s\n```\n\n", cmd.Config.InstallURL, inst)
		installers = append(installers, instText)
	}

	// Replace %PARTIAL_INSTALL% with installers
	partialInstall := strings.TrimSpace(strings.Join(cmd.Sorted(installers), ""))
	markdown = strings.ReplaceAll(markdown, "%PARTIAL_INSTALL%", partialInstall)

	// Tempoarary install explanation
	var tmpFlags string
	for i, name := range installerNames {
		switch i {
		case 0:
			tmpFlags += fmt.Sprintf("`--%s`", name)
		case len(installerNames) - 1:
			tmpFlags += fmt.Sprintf(", or `--%s`", name)
		default:
			tmpFlags += fmt.Sprintf(", `--%s`", name)
		}
	}

	// Replace %TMP_FLAGS% with tmpFlags
	markdown = strings.ReplaceAll(markdown, "%TMP_FLAGS%", tmpFlags)

	// Replace %TMP_DIR% with tmp_dir
	tmp_dir := strings.ReplaceAll(cmd.Config.TmpDir, "@repo_name", cmd.Config.Metadata.Repo)
	markdown = strings.ReplaceAll(markdown, "%TMP_DIR%", tmp_dir)

	// Write markdown to INSTALL.md
	writeString("INSTALL.md", markdown)
}

// Writes README.md and INSTALL.md
func main() {
	writeREADME()
	writeINSTALL()
}
