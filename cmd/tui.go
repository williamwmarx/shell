package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styes for the TUI
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

// Store shell command to run and the name to display when running it
type action struct {
	command string
	name    string
}

// List item
type item string

// FilterValue is required by the list.Item interface.
func (i item) FilterValue() string { return "" }

// itemDelegate is required by the list.Model interface.
type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	// Render the list item
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

// Bubble Tea model
type model struct {
	list             list.Model
	actions          []action
	index            int
	spinner          spinner.Model
	firstFlagInstall bool
	done             bool
	quitting         bool
}

// Initialize the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update the model when a message is received
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// If actions are present, run them
	if len(m.actions) > 1 {
		if m.firstFlagInstall {
			m.firstFlagInstall = false
			return m, tea.Batch(runAction(m.actions[m.index]), m.spinner.Tick)
		}
		return updateChosen(msg, m)
	}
	// Otherwise, show the list of choices
	return updateChoices(msg, m)
}

// Select a choice from the list and add corresponding actions to the queue
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
				case "tmux config":
					m.actions = append(m.actions, tmuxConfig(false)...)
				case "Temporary Zsh config (no plugins)":
					m.actions = append(m.actions, vanillaZshConfig(true)...)
				case "Temporary Vim config (no plugins)":
					m.actions = append(m.actions, vanillaVimConfig(true)...)
				case "Temporary tmux config":
					m.actions = append(m.actions, tmuxConfig(true)...)
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

// Run the next action in the queue
func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case completedActionsMsg:
		m.index++
		if m.index >= len(m.actions) {
			m.done = true
			return m, tea.Quit
		}
		return m, tea.Batch(
			tea.Printf("%s %s", checkMark, m.actions[m.index-1].name),
			runAction(m.actions[m.index]),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// Change view based on model status
func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Cancelling configuration 😔")
	}
	if len(m.actions) > 1 {
		return chosenView(m)
	}
	return choicesView(m)
}

// View for the list of choices
func choicesView(m model) string {
	return "\n" + m.list.View()
}

// View for the current action
func chosenView(m model) string {
	if m.done {
		lastPackageComplete := fmt.Sprintf("%s %s\n", checkMark, m.actions[m.index-1].name)
		return lastPackageComplete + quitTextStyle.Render("All tasks complete 😊")
	}
	info := currentActionStyle.Render(m.actions[m.index].name)
	return fmt.Sprintf("%s%s (%d/%d)", m.spinner.View(), info, m.index+1, len(m.actions))
}

type completedActionsMsg string

// Run a command and return a message when it's done
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

// Run the TUI
func tui(tuiOptions []string) {
	tuiOptions = append([]string{""}, tuiOptions...)

	// Spinner style
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Bold(true)

	// List of actions
	actions := []action{}

	if len(tuiOptions) > 1 {
		// Options passed, run without TUI list selector
		if contains(tuiOptions, "tmux") {
			actions = append(actions, tmuxConfig(false)...)
		}
		if contains(tuiOptions, "tmux temporary") {
			actions = append(actions, tmuxConfig(true)...)
		}
		if contains(tuiOptions, "vim") {
			actions = append(actions, vimConfig()...)
		}
		if contains(tuiOptions, "vanilla-vim") {
			actions = append(actions, vanillaVimConfig(false)...)
		}
		if contains(tuiOptions, "vanilla-vim temporary") {
			actions = append(actions, vanillaVimConfig(true)...)
		}
		if contains(tuiOptions, "zsh") {
			actions = append(actions, zshConfig()...)
		}
		if contains(tuiOptions, "vanilla-zsh") {
			actions = append(actions, vanillaZshConfig(false)...)
		}
		if contains(tuiOptions, "vanilla-zsh temporary") {
			actions = append(actions, vanillaZshConfig(true)...)
		}
		if contains(tuiOptions, "full") {
			actions = fullConfig() // Full config bypasses other options
		}
	}

	// No options passed, launch the TUI list selector
	items := []list.Item{
		item("Full shell config"),
		item("Zsh/Oh My Zsh config"),
		item("Vim + plugins config"),
		item("tmux config"),
		item("Temporary Zsh config (no plugins)"),
		item("Temporary Vim config (no plugins)"),
		item("Temporary tmux config"),
		item("Core packages"),
		item("Design packages"),
		item("Core GUI packages"),
		item("Design GUI packages"),
	}

	// Remove duplicate actions with name "Updating package manager", "Creating .shell.tmp directory" combine uninstall scripts
	filteredActions := []action{}
	alreadyUpdated := false
	alreadyCreatedTmpDir := false
	uninstallCommands := []string{}
	for _, a := range actions {
		if a.name == "Updating package manager" {
			if alreadyUpdated {
				continue
			} else {
				alreadyUpdated = true
			}
		} else if a.name == "Creating .shell.tmp directory" {
			if alreadyCreatedTmpDir {
				continue
			} else {
				alreadyCreatedTmpDir = true
			}
		} else if a.name == "Saving uninstall script" {
			for _, c := range strings.Split(a.command, " && ") {
				tc := strings.TrimSpace(c)
				if tc != "rm -rf ~/.shell.tmp" {
					uninstallCommands = append(uninstallCommands, tc)
				}
			}
			continue
		}
		filteredActions = append(filteredActions, a)
	}
	if len(uninstallCommands) > 0 {
		uninstallCommand := fmt.Sprintf("echo \"%s && rm -rf ~/.shell.tmp\" > ~/.shell.tmp/uninstall.sh", strings.Join(uninstallCommands, " && "))
		filteredActions = append(filteredActions, action{uninstallCommand, "Saving uninstall script"})
	}

	// Setup list
	l := list.New(items, itemDelegate{}, 25, len(items)+6)
	l.Title = "Hi 👋 Let's set up your shell"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// Setup model
	m := model{list: l, spinner: s, actions: filteredActions, firstFlagInstall: len(filteredActions) > 1}

	// Run the program
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
