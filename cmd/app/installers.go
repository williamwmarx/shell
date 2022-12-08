package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Download a file from this repo to a given path
func download(url_path, write_path string) {
	// Get request
	url = "https://raw.githubusercontent.com/williamwmarx/shell/main/" + url_path
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	// Read text from body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Write file to disk
	ioutil.WriteFile(write_path, []byte(body), 0644)
}

// Command exists
func commandExists(commandName string) bool {
	_, err := exec.LookPath(commandName)
	return err == nil
}

// Struct of different pacakge names across installers
type packageName struct {
	aptName string,
	brewName string,
	dnfName string
}

// Install a package
func install(pkg packageName) {
	var cmd string;
	if commandExists("apt-get") {
		cmd = fmt.Sprintf("apt-get install %s", pkg.aptName)
	} else if commandExists("brew") 
		cmd = fmt.Sprintf("brew install %s", pkg.brewName)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Install tmux
func Tmux() {
	fmt.Println("Installing tmux")

	var tmux packageName = {
		aptName: "tmux",
		brewName: "tmux"
		dnfName: "tmux"
	}
	install(tmux)

	download("tmux/tmux.conf", "~/.tmux.conf")
}

func TmpTmux() {
	fmt.Println("Installing tmux in ~/.shell.tmp")
	download("tmux/tmux.conf", "~/.shell.tmp/.tmux.conf")
}

func Vim() {
	fmt.Println("Installing Vim")
	download("vim/vimrc", "~/.vimrc")
}

func TmpVim() {
	fmt.Println("Installing Vim in ~/.shell.tmp")
	download("vim/vimrc", "~/.shell.tmp/.vimrc")
}

func Zsh() {
	fmt.Println("Installing Zsh")
	download("zsh/zshrc", "~/.zshrc")
	download("zsh/aliases", "~/.aliases")
	download("zsh/functions", "~/.functions")
}

func TmpZsh() {
	fmt.Println("Installing Zsh in ~/.shell.tmp")
	download("zsh/zshrc", "~/.shell.tmp/.zshrc")
	download("zsh/aliases", "~/.shell.tmp/.aliases")
	download("zsh/functions", "~/.shell.tmp/.functions")
}
