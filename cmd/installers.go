package cmd

import (
	"fmt"
	"time"
)

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
