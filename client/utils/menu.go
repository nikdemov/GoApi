package utils

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
func WebMenu(inpit string) string {
	p := tea.NewProgram(initialModel(inpit))

	tsk, err := p.StartReturningModel()
	if err != nil {
		log.Fatal(err)
	}

	str := tsk.View()
	split := strings.Split(str, ">")
	str = strings.Join(split, " ")
	split = strings.Split(str, " ")
	str = strings.Join(split, "")
	str = normalform8(str)
	//fmt.Print("Adress:", str)
	return str
}

func normalform8(s string) string {
	if last := len(s) - 8; last >= 0 {
		s = s[:last]
	}
	return s
}

type tickMsg struct{}
type errMsg error

type ipInput struct {
	textInput textinput.Model
	err       error
}

type config struct {
	Ip   []string
	Port string
}

func initialModel(input string) ipInput {
	ti := textinput.NewModel()
	if input == "0" {
		ti.Placeholder = "192.168.0.1"
	} else {
		ti.Placeholder = input

	}

	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ipInput{
		textInput: ti,
		err:       nil,
	}
}

func (m ipInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m ipInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m ipInput) View() string {
	return fmt.Sprintf(

		m.textInput.View())

}

var clear map[string]func() //create a map for storing clear funcs

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		fmt.Println("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
