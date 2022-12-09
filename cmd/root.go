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
	// Add actions
	actions := []action{}
	// Install tmux?
	if options.tmux {
		if options.tmp {
			fmt.Println("Temporarily install tmux")
		} else {
			fmt.Println("Install tmux")
		}
	}
	// Install vim?
	if options.vim {
		if options.tmp {
			fmt.Println("Temporarily install vim")
		} else {
			fmt.Println("Install vim")
		}
	}
	// Install zsh?
	if options.zsh {
		actions = append(actions, action{"zsh", blankFunc})
	}
	// Launch TUI
	tui(actions, options.tmp)
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
