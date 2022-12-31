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
	msg     string
	command string
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

// Handle user input when the list of choices is shown
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Get the selected item and add the corresponding actions to the queue
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if string(i) == "Full shell config" {
					m.actions = append(m.actions, fullConfig()...)
				} else {
					// Iterate through installers to find a match and add the corresponding actions
					for flag, v := range Config.Installers {
						// Get temporary install message
						hm := strings.Fields(v.HelpMessage)
						tmpItemMsg := "Temporarily " + strings.ToLower(hm[0]) + " " + strings.Join(hm[1:], " ")

						if string(i) == v.HelpMessage {
							// Normal install
							m.actions = append(m.actions, install(flag, false)...)
						} else if string(i) == tmpItemMsg {
							// Temporary install
							m.actions = append(m.actions, install(flag, true)...)
						}
					}
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
			tea.Printf(fmt.Sprintf("%s %s", checkMark, m.actions[m.index-1].msg)),
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
		lastPackageComplete := fmt.Sprintf("%s %s\n", checkMark, m.actions[m.index-1].msg)
		return lastPackageComplete + quitTextStyle.Render("All tasks complete 😊")
	}
	info := currentActionStyle.Render(m.actions[m.index].msg)
	return fmt.Sprintf("%s%s (%d/%d)", m.spinner.View(), info, m.index+1, len(m.actions))
}

type completedActionsMsg string

// Run a command and return a message when it's done
func runAction(a action) tea.Cmd {
	return tea.Tick(time.Millisecond*0, func(t time.Time) tea.Msg {
		runCommand(a.command)
		return completedActionsMsg(a.msg)
	})
}

// Run the TUI
func tui(tuiOptions map[string]bool) {
	// Spinner style
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Bold(true)

	// List of actions
	actions := []action{}

	// Parse flags
	if tuiOptions["full"] {
		// Full config
		actions = append(actions, fullConfig()...)
	} else {
		tmp := tuiOptions["tmp"]
		for flag, present := range tuiOptions {
			// Ignore tmp flag, as we already recorded its value
			if flag == "tmp" {
				continue
			}
			// If flag is present, add the corresponding actions
			if present {
				actions = append(actions, install(flag, tmp)...)
			}
		}

		// Remove duplicate actions
		exportActions := []action{}
		executedCommands := []string{}
		for _, a := range actions {
			if !contains(executedCommands, a.command) {
				exportActions = append(exportActions, a)
				executedCommands = append(executedCommands, a.command)
			}
		}
		actions = exportActions
	}

	// No options passed, launch the TUI list selector
	items := []list.Item{item("Full shell config")}

	// Add installers to the list
	for _, v := range Config.Installers {
		items = append(items, item(v.HelpMessage))
	}

	// Add temporary installers to the list
	for _, v := range Config.Installers {
		// Create temporary help message and append to items
		hm := strings.Fields(v.HelpMessage)
		message := "Temporarily " + strings.ToLower(hm[0]) + " " + strings.Join(hm[1:], " ")
		items = append(items, item(message))
	}

	// Add package groups to the list
	for packageGroup := range PM.packages {
		items = append(items, item(packageGroup+" packages"))
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
	m := model{list: l, spinner: s, actions: actions, firstFlagInstall: len(actions) > 1}

	// Run the program
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
