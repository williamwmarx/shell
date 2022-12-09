package cmd

import (
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"os"
)

// Potential install options
type Options struct {
	tmp bool
	tmux bool
	vim bool
	zsh bool
}

func runRoot(options Options) {
	// Track whether any flags found
	flagsPresent := false
	// Install tmux?
	if options.tmux {
		flagsPresent = true
		if options.tmp {
			fmt.Println("Temporarily install tmux")
		} else {
			fmt.Println("Install tmux")
		}
	}
	// Install vim?
	if options.vim {
		flagsPresent = true
		if options.tmp {
			fmt.Println("Temporarily install vim")
		} else {
			fmt.Println("Install vim")
		}
	}
	// Install vim?
	if options.zsh {
		flagsPresent = true
			fmt.Println("Temporarily install zsh")
		if options.tmp {
		} else {
			fmt.Println("Install zsh")
		}
	}
	// No flags found, launch TUI
	if !flagsPresent {
		tui()
	}
}

func flagPresent(cmd *cobra.Command, flagName string) bool {
	isPresent, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		log.Fatal(err)
	}
	return isPresent
}

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Install my default packages and dotfiles",
	Run:   func(cmd *cobra.Command, args []string) {
		options := Options{
			tmp: flagPresent(cmd, "tmp"),
			tmux: flagPresent(cmd, "tmux"),
			vim: flagPresent(cmd, "vim"),
			zsh: flagPresent(cmd, "zsh"),
		}
		runRoot(options)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("tmp", "", false, "Only temporarily install selection(s)")
	rootCmd.Flags().BoolP("tmux", "", false, "Install tmux and configuration files")
	rootCmd.Flags().BoolP("vim", "", false, "Install vim and configuration files")
	rootCmd.Flags().BoolP("zsh", "", false, "Install zsh and configuration files")
}
