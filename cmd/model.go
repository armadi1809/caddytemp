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
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
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
	list             list.Model
	choice           string
	quitting         bool
	step             int
	selectedTemplate string
	inputs           []textinput.Model
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
				}
				return m, nil
			case 1:
				return m, m.updateInputs(msg)
				// fileContent, ok := templates[m.selectedTemplate]
				// if !ok {
				// 	m.choice = "An unexpected error occurred."
				// 	return m, tea.Quit
				// }
				// f, err := os.Create("Caddyfile")
				// if err != nil {
				// 	m.choice = "An unexpected error occurred."
				// 	return m, tea.Quit
				// }
				// defer f.Close()
				// userInputs := []any{}
				// for _, placeholder := range placeholdersData[m.selectedTemplate] {
				// 	fmt.Printf("Enter %s: ", placeholder)
				// 	reader := bufio.NewReader(os.Stdin)
				// 	line, err := reader.ReadString('\n')
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// 	userInputs = append(userInputs, line)
				// }
				// _, err = f.WriteString(fmt.Sprintf(fileContent, userInputs...))
				// if err != nil {
				// 	m.choice = "An unexpected error occurred."
				// 	return m, tea.Quit
				// }
				// f.Sync()
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
	return b.String()
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
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, step: 0}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
