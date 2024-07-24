package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle       = focusedStyle
	noStyle           = lipgloss.NewStyle()
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list                   list.Model
	choice                 string
	quitting               bool
	step                   int
	selectedTemplate       string
	inputs                 []textinput.Model
	selectedPlaceholderIdx int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch m.step {
			case 0:
				m.initializeInputs()
				return m, nil
			case 1:
				if (m.selectedPlaceholderIdx) == len(placeholdersData[m.selectedTemplate])-1 {
					err := m.writeInputsToFile()
					if err != nil {
						m.choice = err.Error()
					}
					return m, tea.Quit
				}
				m.selectedPlaceholderIdx += 1
				cmds := m.updateModelInputs()
				return m, tea.Batch(cmds...)
			}
			return m, tea.Quit
		}
	}
	if m.step == 1 {
		return m, m.updateInputs(msg)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
func (m model) View() string {
	if m.step == 0 {
		return "\n" + m.list.View()
	}
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Generating Caddyfile: %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Quitting...")
	}
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.list.View(), b.String())
}

func (m *model) updateModelInputs() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.selectedPlaceholderIdx {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}
	return cmds
}

func (m *model) writeInputsToFile() error {
	fileContent, ok := templates[m.selectedTemplate]
	if !ok {
		return fmt.Errorf("an unexpected error occurred")
	}
	f, err := os.Create("Caddyfile")
	if err != nil {
		return fmt.Errorf("an unexpected error occurred")
	}
	defer f.Close()
	userInputs := map[string]string{}
	for _, inp := range m.inputs {
		userInputs[inp.Placeholder] = inp.Value()
	}
	for _, pl := range placeholdersData[m.selectedTemplate] {
		fileContent = strings.ReplaceAll(fileContent, "{"+pl+"}", userInputs[pl])
	}
	_, err = f.WriteString(fileContent)
	if err != nil {
		return fmt.Errorf("an unexpected error occurred")
	}
	f.Sync()
	return nil
}

func (m *model) initializeInputs() {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		m.selectedTemplate = string(i)
		placeholders := placeholdersData[string(i)]
		m.inputs = make([]textinput.Model, len(placeholders))
		var t textinput.Model
		for i := range m.inputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			t.CharLimit = 32
			t.Placeholder = placeholders[i]
			if i == 0 {
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
			}
			m.inputs[i] = t
		}
		m.step += 1
		m.selectedPlaceholderIdx = 0
	}
}

func run() {
	items := []list.Item{}
	for k := range templates {
		items = append(items, item(k))
	}
	const defaultWidth = 20
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select a Cadddyfile template"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, step: 0}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
