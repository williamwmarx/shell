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

func main() {
	writePackagesREADME()
}
