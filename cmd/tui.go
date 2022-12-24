package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle         = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("170"))
	itemStyle          = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle          = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle      = lipgloss.NewStyle().Margin(1, 0, 2, 4).Bold(true)
	currentActionStyle = lipgloss.NewStyle().Bold(true)
	checkMark          = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

type action struct {
	command string
	name    string
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list             list.Model
	actions          []action
	index            int
	spinner          spinner.Model
	firstFlagInstall bool
	done             bool
	quitting         bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	if len(m.actions) > 1 {
		if m.firstFlagInstall {
			m.firstFlagInstall = false
			return m, tea.Batch(runAction(m.actions[m.index]), m.spinner.Tick)
		}
		return updateChosen(msg, m)
	}
	return updateChoices(msg, m)
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				switch string(i) {
				case "Full shell config":
					m.actions = append(m.actions, fullConfig()...)
				case "Zsh/Oh My Zsh config":
					m.actions = append(m.actions, zshConfig()...)
				case "Vim + plugins config":
					m.actions = append(m.actions, vimConfig()...)
				case "[TMP] Zsh config (no plugins)":
					m.actions = append(m.actions, action{"", "[TMP] Zsh config (no plugins)"})
				case "[TMP] Vim config (no plugins)":
					m.actions = append(m.actions, action{"", "[TMP] Vim config (no plugins)"})
				case "Core packages":
					m.actions = append(m.actions, installActions("Core")...)
				case "Design packages":
					m.actions = append(m.actions, installActions("Design")...)
				case "Core GUI packages":
					m.actions = append(m.actions, installActions("GuiCore")...)
				case "Design GUI packages":
					m.actions = append(m.actions, installActions("GuiDesign")...)
				}
				return m, tea.Batch(runAction(m.actions[m.index]), m.spinner.Tick)
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case completedActionsMsg:
		if m.index >= len(m.actions)-1 {
			m.done = true
			return m, tea.Quit
		}

		m.index++
		return m, tea.Batch(
			tea.Printf("%s %s", checkMark, m.actions[m.index].name),
			runAction(m.actions[m.index]),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Cancelling configuration 😔")
	}
	if len(m.actions) > 1 {
		return chosenView(m)
	}
	return choicesView(m)
}

func choicesView(m model) string {
	return "\n" + m.list.View()
}

func chosenView(m model) string {
	if m.done {
		return quitTextStyle.Render("All tasks complete 😊")
	}
	info := currentActionStyle.Render(m.actions[m.index].name)
	return fmt.Sprintf("%s%s (%d/%d)", m.spinner.View(), info, m.index, len(m.actions)-1)
}

type completedActionsMsg string

func runAction(a action) tea.Cmd {
	return tea.Tick(time.Millisecond*0, func(t time.Time) tea.Msg {
		runCommand(a.command)
		return completedActionsMsg(a.name)
	})
}

// Function to check if array of string contains element
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func tui(tuiOptions []string) {
	tuiOptions = append([]string{""}, tuiOptions...)

	// Spinner style
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Bold(true)

	// List of actions
	actions := []action{{"echo -n", ""}}

	if len(tuiOptions) > 1 {
		// Options passed, run without TUI list selector
		if contains(tuiOptions, "tmux") {
			actions = append(actions, tmuxConfig()...)
		} else if contains(tuiOptions, "tmux temporary") {
			// TODO: tmux temporary
		} else if contains(tuiOptions, "vim") {
			actions = append(actions, vimConfig()...)
		} else if contains(tuiOptions, "vanilla-vim") {
			// TODO: vanilla-vim
		} else if contains(tuiOptions, "vanilla-vim temporary") {
			// TODO: vanilla-vim temporary
		} else if contains(tuiOptions, "zsh") {
			actions = append(actions, zshConfig()...)
		} else if contains(tuiOptions, "vanilla-zsh") {
			// TODO: vanilla-zsh
		} else if contains(tuiOptions, "vanilla-zsh temporary") {
			// TODO: vanilla-zsh temporary
		} else {
			log.Fatal("Invalid option")
			os.Exit(1)
		}
	}

	// No options passed, launch the TUI list selector
	items := []list.Item{
		item("Full shell config"),
		item("Zsh/Oh My Zsh config"),
		item("Vim + plugins config"),
		item("tmux config"),
		item("[TMP] Zsh config (no plugins)"),
		item("[TMP] Vim config (no plugins)"),
		item("Core packages"),
		item("Design packages"),
		item("Core GUI packages"),
		item("Design GUI packages"),
	}

	l := list.New(items, itemDelegate{}, 25, len(items)+6)
	l.Title = "Hi 👋 Let's set up your shell"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, spinner: s, actions: actions, firstFlagInstall: len(actions) > 1}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
