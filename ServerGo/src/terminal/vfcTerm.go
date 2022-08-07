package terminal

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"os"

	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Vmodel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() Vmodel {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Vmodel{spinner: s}
}

func (m Vmodel) Init() tea.Cmd {
	return spinner.Tick
}

func (m Vmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m Vmodel) View() string {
	if m.err != nil {
		logs.ErrorLogger.Println("VFC Screen" + m.err.Error())
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s VFS Service running ... press q to quit\n\n and press ctrl+c to stop program\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}

func VFSTerm() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
