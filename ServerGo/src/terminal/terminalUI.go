package terminal

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	generator "nikworkedprofile/GoApi/ServerGo/src/generate_logs"
	"nikworkedprofile/GoApi/ServerGo/src/logenc"
	"nikworkedprofile/GoApi/ServerGo/src/web"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

// General stuff for styling the view
var (
	pathdata      = "/var/local/logi2"
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)
)

type tickMsg struct{}
type frameMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type Model struct {
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

var (
	timeout0 = Model{0, false, 0, 0, 0, false, true}
	timeout1 = Model{1, false, 0, 0, 0, false, true}
	timeout2 = Model{2, false, 0, 0, 0, false, true}
	timeout3 = Model{3, false, 0, 0, 0, false, true}
	timeout4 = Model{4, false, 0, 0, 0, false, true}
	timeout5 = Model{5, false, 0, 0, 0, false, true}

	status tea.Model
	//test   []string
	//ctx, _ = context.WithCancel(context.Background())
)

func MainUi() {
	var test tea.Model
	var st bool
	str, model := TerminalUi()
	idx, _ := strconv.Atoi(str)
	if (model == timeout0 || model == timeout1 || model == timeout2 || model == timeout3 || model == timeout4 || model == timeout5) && idx == 0 {
		test = Screensaver()
	} else if model != timeout0 || model != timeout1 || model != timeout2 || model != timeout3 || model != timeout4 || model != timeout5 {
		st = SwitchMenu(idx)

	}
	if test != nil || status != nil || st {
		status = nil
		test = nil
		MainUi()
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	var s string
	if m.Quitting {
		//return m, tea.Quit
		str := strconv.Itoa(m.Choice)
		return str
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice += 1
			if m.Choice > 8 {
				m.Choice = 8
			}
		case "k", "up":
			m.Choice -= 1
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks -= 1
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	//var s string
	switch msg.(type) {

	case frameMsg:
		if !m.Loaded {
			m.Frames += 1
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Ticks = 5
				return m, tick()
			}
			return m, frame()
		}

	case tickMsg:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				//View()
				return m, tea.Quit
			}
			m.Ticks -= 1
			return m, tick()
		}
	}

	return m, nil
}

// Sub-views

// The first view, where you're choosing a task
func choicesView(m Model) string {
	c := m.Choice

	tpl := "Control panel\n\n"
	tpl += "%s\n\n"
	tpl += "Program in wait mode %s seconds\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s",
		checkbox("Gen logs", c == 0),
		checkbox("Run Web", c == 1),
		checkbox("running VFS", c == 2),
		checkbox("clear genlogs", c == 3),
		checkbox("Search word or collocation", c == 4),
		checkbox("Enter Ip server for collect logs", c == 5),
	)

	return fmt.Sprintf(tpl, choices, colorFg(strconv.Itoa(m.Ticks), "79"))
}

// The second view, after a task has been chosen
func chosenView(m Model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("GenLogs\n\nCool, we generate logs %s and %s...", keyword("size generate logs"), keyword("Count generate logs"))

	case 1:
		port := "15000"
		msg = fmt.Sprintf("Run Web\n\n Start web interface ...%s.", keyword(port))
	case 2:
		msg = fmt.Sprintf("running VFS\n\n We start VFS service  %s ...", keyword("OK"))
	case 3:
		msg = fmt.Sprintf("clear genlogs\n\n Please wait, we clear generated...")
		generator.RemoveByConfig()
	case 4:
		fmt.Print("Enter content for Search:")
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')
		//logenc.SearchT(text)
		msg = fmt.Sprintf("Search word or collocation\n\nPlease enter word or collocation ")
		//logenc.SearchT(text)
	case 5:
		msg = fmt.Sprintf("EnterIp\n\nCool, running form and data Ip")
		//web.EnterIp()
	default:
		msg = fmt.Sprintf("Okay.\n\nYou enter the error please restart program /n/n Report a bug in %s or %s...", keyword("Contact 1"), keyword("Contact 2"))
	}

	label := "Loading..."
	if m.Loaded {
		label = fmt.Sprintf("Loaded. Following a %s seconds...", colorFg(strconv.Itoa(m.Ticks), "79"))
	}

	return msg + "\n\n" + label + "\n" + progressbar(80, m.Progress) + "%"
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(width int, percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
/* func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
} */

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
func SwitchMenu(idx int) (exit bool) {

	switch choose := idx; choose {
	case 0:
		//TODO: add size file and count file
		//UI for gerator animation when logs generated
		generator.ProcGenN(10, 2000)
		exit = true
	case 1:
		//UI for run web and main server
		//add in config file
		//????
		//fmt.Print("Enter port for run Web:")
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')
		ctxWEB, err := context.WithCancel(context.Background())
		if err != nil {
			log.Print(err)

		}
		var test []string
		CallClear()
		web.ProcWeb("-p", test, ctxWEB)
		exit = true
	case 2:
		//UI for run VFS animation
		//ctx, _ := context.WithCancel(context.Background())
		//go controllers.VFS("10015", ctx)
		VFSTerm()
		//time.Sleep(5 * time.Second)
		exit = true

		//go VFSTerm()
	case 3:
		//UI for example animation
		//add case for clear reddata
		generator.RemoveByConfig()
		exit = true
	case 4:
		//add UI for search in terminal
		//fmt.Print("Enter content for Search:")
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')
		logenc.SearchT(pathdata + "/repdata/")
		exit = true
	case 5:
		web.EnterIp()
		exit = true
	}
	return exit
}

func TerminalUi() (string, tea.Model) {

	initialModel := Model{0, false, 10, 0, 0, false, false}
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	model, err := p.StartReturningModel()
	if err != nil {
		fmt.Println("could not start program:", err)
	}

	str := model.View()

	return str, model
}

/* func VfsUiTerm() {

}
func WebUiTerm() {

}
func ProcFileUiTerm() {

} */

var clear map[string]func() //create a map for storing clear funcs

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
