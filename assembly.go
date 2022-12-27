package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/williamwmarx/shell/cmd"
)

// Store package metadata
type packageMetadata struct {
	name        string
	url				  string
	description string
}

// Writes packages/README.md, a list of all packages
func writePackagesREADME() {
	// Read the file packages README template
	markdownBase, err := os.ReadFile("assets/packages_README_base.md")
	if err != nil {
		log.Fatal(err)
	}
	markdown := string(markdownBase)

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
	f, err := os.Create("packages/README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(strings.TrimSpace(markdown))
}

type targetsText struct {
	header string
	body string
}

// Writes apex README.md, showing install command and listing synced dotfiles
func writeApexREADME() {
	// Get install URL
	installURL := cmd.Config.CustomInstallURL
	if installURL == "" {
		installURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/main/install.sh", cmd.RepoPath())
	}

	// Read the apex README template
	markdownBase, err := os.ReadFile("assets/apex_README_base.md")
	if err != nil {
		log.Fatal(err)
	}
	// Add proper URL for install script
	markdown := strings.ReplaceAll(string(markdownBase), "INSTALL_URL", installURL)

	// Format text for synced targets
	var syncTargets []targetsText
	for _, s := range cmd.Config.Sync {
		sync := reflect.ValueOf(&s).Elem()
		// Name of target group as h3
		header := fmt.Sprintf("### %s\n", sync.FieldByName("Name").String())

		// Description of target group
		body := sync.FieldByName("Description").String()
		if body != "" {
			body = fmt.Sprintf("%s\n", body)
		}

		// List of targets and paths
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
	f, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(strings.TrimSpace(markdown))
}

func main() {
	writePackagesREADME()
	writeApexREADME()
}
