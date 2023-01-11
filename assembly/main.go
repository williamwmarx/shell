package main

import (
	"fmt"
	"log"
	"os"
	"sort"
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
		header += cmd.PM.Packages[packageGroup].Description + "\n"
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
	installers := cmd.Config.Installers
	var installerNames []string
	for k := range installers {
		installerNames = append(installerNames, k)
	}
	sort.Strings(installerNames)

	// Add installers to markdown
	for _, in := range installerNames {
		// Title and description
		markdown += fmt.Sprintf("#### %s\n%s\n", in, installers[in].Description)
		// Code block
		markdown += fmt.Sprintf("```bash\nsh <(curl %s) --%s\n```\n", cmd.Config.InstallURL, in)
		// Extra line break
		markdown += "\n"
	}

	// Tempoarary install explanation
	tmp_explanation := "### Temporary install\nSometimes, you only need your dotfiles " +
		"temporarily. For example, say you're editing some code on a friend's machine. " +
		"You could slowly go through it with their editor, or you could load up your vim " +
		"config and fly through their code. This is where the `--tmp` flag comes in. You " +
		"can use the `--tmp` flag with "
	// Add installers flags
	for i, name := range installerNames {
		switch i {
		case 0:
			tmp_explanation += fmt.Sprintf("`--%s`", name)
		case len(installerNames) - 1:
			tmp_explanation += fmt.Sprintf(", or `--%s`", name)
		default:
			tmp_explanation += fmt.Sprintf(", `--%s`", name)
		}
	}
	tmp_explanation += ". It will install the packages, download necessary dotfiles into " +
		"the `TMP_DIR` directory, and add the shell script `TMP_DIR/uninstall.sh` which " +
		"will uninstall any packages you installed and remove the `TMP_DIR` directory. " +
		"Temporary install will look for the “vanilla” versions of synced dotfiles, where " +
		"possible."
	tmp_dir := strings.ReplaceAll(cmd.Config.TmpDir, "@repo_name", cmd.Config.Metadata.Repo)
	markdown += strings.ReplaceAll(tmp_explanation, "TMP_DIR", tmp_dir)

	// Write markdown to INSTALL.md
	writeString("INSTALL.md", markdown)
}

// Writes README.md and INSTALL.md
func main() {
	writeREADME()
	writeINSTALL()
}
