package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/williamwmarx/shell/cmd"
)

// Store package metadata
type packageMetadata struct {
	name        string
	url         string
	description string
}

// Writes packages/README.md, a list of all packages
func writePackagesREADME() {
	// Read the file packages README template
	markdown := cmd.ReadToString("assets/packages_README_base.md")

	// Iterate through package groups
	packageGroups := reflect.ValueOf(&cmd.Packages).Elem()
	for i := 0; i < packageGroups.NumField(); i++ {
		// Package group header
		markdown += fmt.Sprintf("### %s\n", packageGroups.Type().Field(i).Name)

		// Collect all packages in group
		var packages []packageMetadata
		p := packageGroups.Field(i).MapRange()
		for p.Next() {
			name := p.Value().FieldByName("Name").String()
			url := p.Value().FieldByName("URL").String()
			description := p.Value().FieldByName("Description").String()
			packages = append(packages, packageMetadata{name, url, description})
		}

		// Sort packages by name
		sort.Slice(packages, func(i, j int) bool {
			return strings.ToLower(packages[i].name) < strings.ToLower(packages[j].name)
		})

		// Append package info to markdown
		for _, p := range packages {
			markdown += fmt.Sprintf("- [%s](%s) - %s\n", p.name, p.url, p.description)
		}

		// Add extra line break
		markdown += "\n"
	}

	// Write markdown to packages/README.md
	cmd.WriteString("packages/README.md", markdown)
}

// Store synced target markdown format
type targetsText struct {
	header string
	body   string
}

// Writes apex README.md, showing install command and listing synced dotfiles
func writeApexREADME() {
	// Read the apex README template and insert proper URL
	markdown := cmd.ReadToString("assets/apex_README_base.md")
	markdown = strings.ReplaceAll(markdown, "INSTALL_URL", cmd.Metadata.URL)

	// Format text for synced targets
	var syncTargets []targetsText
	for _, s := range cmd.Config.Sync {
		sync := reflect.ValueOf(&s).Elem()
		// Name of target group as h3
		header := fmt.Sprintf("### %s\n", sync.FieldByName("Name").String())

		// List of targets and paths
		var body string
		for _, t := range sync.FieldByName("Targets").Interface().([]cmd.Target) {
			paths := strings.Split(t.LocalPath, "/")
			body += fmt.Sprintf("- [%s](%s) — %s\n", paths[len(paths)-1], t.RepoPath, t.Description)
		}

		syncTargets = append(syncTargets, targetsText{header, body + "\n"})
	}

	// Sort sync targets by name
	sort.Slice(syncTargets, func(i, j int) bool {
		return strings.ToLower(syncTargets[i].header) < strings.ToLower(syncTargets[j].header)
	})

	// Add sync targets to markdown
	for _, d := range syncTargets {
		markdown += d.header + d.body
	}

	// Write markdown to README.md
	cmd.WriteString("README.md", markdown)
}

// Writes INSTALL.md, showing thorough install options
func writeINSTALL() {
	// Read the apex README template and insert proper URL
	markdown := cmd.ReadToString("assets/INSTALL_base.md")
	markdown = strings.ReplaceAll(markdown, "INSTALL_URL", cmd.Metadata.URL)

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
		markdown += fmt.Sprintf("```bash\nsh <(curl %s) --%s\n```\n", cmd.Metadata.URL, in)
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
	tmp_dir := strings.ReplaceAll(cmd.Config.TmpDir, "@repo_name", cmd.Metadata.Name)
	markdown += strings.ReplaceAll(tmp_explanation, "TMP_DIR", tmp_dir)

	// Write markdown to INSTALL.md
	cmd.WriteString("INSTALL.md", markdown)
}

func main() {
	writePackagesREADME()
	writeApexREADME()
	writeINSTALL()
}
