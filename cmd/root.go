package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Potential install options
type Options struct {
	full       bool
	tmp        bool
	tmux       bool
	vim        bool
	zsh        bool
	vanillaVim bool
	vanillaZsh bool
}

func runRoot(options Options) {
	// Options to pass to TUI
	tuiOptions := []string{}
	// Install tmux?
	if options.tmux {
		if options.tmp {
			tuiOptions = append(tuiOptions, "tmux temporary")
		} else {
			tuiOptions = append(tuiOptions, "tmux")
		}
	}
	// Install vim?
	if options.vim && options.tmp {
		tuiOptions = append(tuiOptions, "vanilla-vim temporary")
	} else if options.vanillaVim {
		tuiOptions = append(tuiOptions, "vanilla-vim")
	} else if options.vim {
		tuiOptions = append(tuiOptions, "vim")
	}
	// Install zsh?
	if options.zsh && options.tmp {
		tuiOptions = append(tuiOptions, "vanilla-zsh temporary")
	} else if options.vanillaZsh {
		tuiOptions = append(tuiOptions, "vanilla-zsh")
	} else if options.zsh {
		tuiOptions = append(tuiOptions, "zsh")
	}
	// Full system config?
	if options.full {
		tuiOptions = []string{"full"}
	}
	tui(tuiOptions)
}

func flagPresent(cmd *cobra.Command, flagName string) bool {
	isPresent, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		log.Fatal(err)
	}
	return isPresent
}

var rootCmd = &cobra.Command{
	Use:   "sh <(curl https://marx.sh)",
	Short: "Install my default packages and dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		options := Options{
			full:       flagPresent(cmd, "full"),
			tmp:        flagPresent(cmd, "tmp"),
			tmux:       flagPresent(cmd, "tmux"),
			vim:        flagPresent(cmd, "vim"),
			vanillaVim: flagPresent(cmd, "vanilla-vim"),
			zsh:        flagPresent(cmd, "zsh"),
			vanillaZsh: flagPresent(cmd, "vanilla-zsh"),
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
	rootCmd.Flags().BoolP("full", "", false, "Full system config")
	rootCmd.Flags().BoolP("tmp", "", false, "Only temporarily install selection(s)")
	rootCmd.Flags().BoolP("tmux", "", false, "Install tmux and configuration files")
	rootCmd.Flags().BoolP("vim", "", false, "Install vim and configuration files")
	rootCmd.Flags().BoolP("vanilla-vim", "", false, "Install vim and configuration files without plugins")
	rootCmd.Flags().BoolP("zsh", "", false, "Install zsh and configuration files")
	rootCmd.Flags().BoolP("vanilla-zsh", "", false, "Install zsh and configuration files without plugins")
}
