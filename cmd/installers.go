package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

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

func decode() {
	var packages pkgs
	_, err := toml.DecodeFile("packages/packages.toml", &packages)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range packages.Core {
		fmt.Println(p.Name)
	}
}

func installZsh(temporary bool) {
	time.Sleep(time.Millisecond * 500)
	fmt.Println("------------- INSTALLING ZSH -----------------")
}

func tmux(temporary bool) {
	time.Sleep(time.Millisecond * 500)
	if temporary {
		fmt.Println("TEMPORARY TMUX")
	} else {
		fmt.Println("PERMANENT TMUX")
	}
}

func vim(temporary bool) {
	time.Sleep(time.Millisecond * 500)
	if temporary {
		fmt.Println("TEMPORARY VIM")
	} else {
		fmt.Println("PERMANENT VIM")
	}
}
