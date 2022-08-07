package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	ipall             []string
)

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

	fmt.Fprintf(w, fn(str))
}

type model struct {
	list     list.Model
	items    []item
	choice   string
	quitting bool
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
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {

	if m.choice != "" {
		Ready(m.choice)

		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

func Edit() {
	//var port string
	var ip []string

	ip, _ = CheckConfig()
	if len(ip) == 0 {
		ip = Ready("0")

	}
	ipall = ip
	items := []list.Item{}
	//items := []list.Item{}
	for i := 0; i < len(ip); i++ {
		items = append(items, item(ip[i]))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func EnterIp(Ipready []string) {
	var data []byte
	var ipaddr []string
	fmt.Print(Ipready)

	for i := 0; i < len(Ipready); i++ {

		ipaddr = append(ipaddr, Ipready[i])
		ipaddr = removeDuplicateStr(ipaddr)
		config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "15000"}}

		data, _ = toml.Marshal(&config)
	}

	err3 := ioutil.WriteFile(pathdata+"/config.toml", data, 0666)

	if err3 != nil {

		log.Fatal(err3)
	}
	fmt.Println(ipaddr)
	fmt.Println("Written")

}

func Ready(choice string) []string {
	var ipc []string
	ipc = ipall
	fmt.Println(ipc)
	fmt.Println(ipall)
	for {
		text := WebMenu(choice)
		if text == "stop" {
			break
		}
		time.Sleep(1e9)

		if CheckIPAddress(text) || text == "0" {
			if choice != "0" {
				for i := 0; i < len(ipc); i++ {

					if choice == ipc[i] {
						ipc[i] = ipc[len(ipc)-1] // Copy last element to index i.
						ipc[len(ipc)-1] = ""     // Erase last element (write zero value).
						ipc = ipc[:len(ipc)-1]
						break
					}
					fmt.Println("ipc", ipc)

				}
				if text != "0" {
					ipc = append(ipc, text)
					fmt.Println(ipc)
					EnterIp(ipc)
				} else {
					EnterIp(ipc)
					fmt.Println("Remove")
					fmt.Println("Enter Ip or send srop")
				}
			}
		}
	}

	return ipc
}

func CheckIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
		return false
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	}

}
