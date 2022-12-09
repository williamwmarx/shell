package cmd

import (
	"fmt"
	"io"
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

func blankFunc(b bool) {}

type action struct {
	name string
	fn   func(bool)
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
	list     list.Model
	actions  []action
	index    int
	spinner  spinner.Model
	done     bool
	quitting bool
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

	if len(m.actions) == 1 {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
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
					m.actions = append(m.actions, action{"Full shell config", blankFunc})
				case "Zsh/Oh My Zsh config":
					m.actions = append(m.actions, action{"Zsh/Oh My Zsh config", installZsh})
				case "Vim + plugins config":
					m.actions = append(m.actions, action{"Vim + plugins config", blankFunc})
				case "[TMP] Zsh config (no plugins)":
					m.actions = append(m.actions, action{"[TMP] Zsh config (no plugins)", blankFunc})
				case "[TMP] Vim config (no plugins)":
					m.actions = append(m.actions, action{"[TMP] Vim config (no plugins)", blankFunc})
				}
				return m, tea.Batch(downloadAndInstall(m.actions[m.index]), m.spinner.Tick)
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
	case installedPkgMsg:
		if m.index >= len(m.actions)-1 {
			m.done = true
			return m, tea.Quit
		}

		m.index++
		return m, tea.Batch(
			tea.Printf("%s %s", checkMark, m.actions[m.index].name),
			downloadAndInstall(m.actions[m.index]),
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
	if len(m.actions) == 1 {
		return choicesView(m)
	}
	return chosenView(m)
}

func choicesView(m model) string {
	return "\n" + m.list.View()
}

func chosenView(m model) string {
	if m.done {
		doneMsg := fmt.Sprintf("Done! Installed %d actions.", len(m.actions)-1)
		return quitTextStyle.Render(doneMsg)
	}

	info := currentActionStyle.Render("Installing " + m.actions[m.index].name)

	return fmt.Sprintf("%s%s (%d/%d)", m.spinner.View(), info, m.index, len(m.actions)-1)
	
}

type installedPkgMsg string

func downloadAndInstall(a action) tea.Cmd {
	return tea.Tick(time.Millisecond*0, func(t time.Time) tea.Msg {
		a.fn(false)
		return installedPkgMsg(a.name)
	})
}

func tui_list() {
	items := []list.Item{
		item("Full shell config"),
		item("Zsh/Oh My Zsh config"),
		item("Vim + plugins config"),
		item("tmux config"),
		item("[TMP] Zsh config (no plugins)"),
		item("[TMP] Vim config (no plugins)"),
	}

	l := list.New(items, itemDelegate{}, 25, len(items)+6)
	l.Title = "Hi 👋 Let's set up your shell"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Bold(true)

	m := model{list: l, spinner: s, actions: []action{action{"", blankFunc}}}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
