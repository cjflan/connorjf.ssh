package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	page      page
	term      string
	profile   string
	width     int
	height    int
	bg        string
	err       error
	viewport viewport.Model
	navStyle lipgloss.Style
	errorStyle lipgloss.Style
	headerStyle lipgloss.Style
	titleStyle lipgloss.Style
	paraStyle lipgloss.Style
	txtStyle  lipgloss.Style
	listStyle lipgloss.Style
	quitStyle lipgloss.Style
}

type page int
const (
	home page = iota
	blog
	resume
	about
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 10
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h":
			m.page = home
		case "b":
			m.page = blog
		case "r":
			m.page = resume
		case "a":
			m.page = about
		}
		if m.page == resume {
			switch msg.String() {
			case "j":
				m.viewport.LineDown(2)
			case "k":
				m.viewport.LineUp(2)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	var page strings.Builder
	page.WriteString(m.headerView())

	switch m.page {
	case home:
		page.WriteString(m.homeView())
	case blog:
		page.WriteString(m.blogView())
	case resume:
		page.WriteString(m.resumeView())
	case about:
		page.WriteString(m.aboutView())
	}
	page.WriteString("\n\n")

	if m.err != nil {
		e := fmt.Sprintf("Error with application: %s", m.err)
		page.WriteString(m.errorStyle.Render(e))
		page.WriteString("\n\n")
	}

	page.WriteString(m.quitStyle.Render("Press 'q' to quit\n"))

	return page.String()
}

func (m model) headerView() string {
	pages := []string{"h home", "b blog", "r resume", "a about"}
	for i, page := range pages {
		if i == int(m.page) {
			pages[i] = m.headerStyle.Render(page)
		} else {
			pages[i] = m.headerStyle.Render(string(page[0])) + m.navStyle.Render(page[1:])
		}
	}

	h := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(m.navStyle).
		Row(pages...)
	
	return h.String() + "\n\n"
}

func (m model) homeView() string {
	var s strings.Builder
	s.WriteString(m.titleStyle.Render("who am i?"))
	s.WriteString("\n\n")
	s.WriteString(m.paraStyle.Render("connor offline, connorjf/enso online"))
	s.WriteString("\n\n")


	s.WriteString(m.titleStyle.Render("how would i describe myself?"))
	s.WriteString("\n\n")
	s.WriteString(m.paraStyle.Render("engineer. tech enthusiast. outdoors lover."))
	s.WriteString("\n\n")


	s.WriteString(m.titleStyle.Render("how've i gotten here?"))
	s.WriteString("\n\n")
	s.WriteString(m.paraStyle.Render("i went to college at the Georgia Institute of Technology for a BS in Aerospace Engineering. These days i'm  working at MacStadium leading the Sales Engineering team. When I’m not at work I can be found going to EDM concerts, rock climbing, sailing, or trying to figure out when I’m next going skiing."))
	s.WriteString("\n\n")
	s.WriteString(m.paraStyle.Render("i am open to a new position in the software engineering space. you can find my resume by pressing 'r'!"))
	s.WriteString("\n\n")

	s.WriteString(m.titleStyle.Render("what am i up to now?"))
	s.WriteString("\n\n")

	s.WriteString(m.paraStyle.Foreground(lipgloss.Color("#0db9d7")).Render("things i'm working on"))
	s.WriteString("\n")
	todo := []string{
		"send a v10 boulder and climb 5.13 (in the gym)", 
		"learn to trad climb", 
		"train to run a marathon (goal sub 3:30)",
		"bike a century",
		"cross country ski (classic) the american birkibeiner",
	}

	t := m.paraStyle.Render(list.New(todo).String())
	s.WriteString(t)
	s.WriteString("\n\n")

	s.WriteString(m.paraStyle.Foreground(lipgloss.Color("#0db9d7")).Render("things i want to do longer term"))
	s.WriteString("\n")
	goals := []string{
		"complete an ironman",
		"hike the pct",
		"learn to backcountry ski",
		"learn to sail a laser",
		"bikepack the iceland ring road",
	}
	g := m.paraStyle.Render(list.New(goals).String())
	s.WriteString(g)
	return s.String()
}

func (m model) blogView() string {
	var s strings.Builder
	s.WriteString("this is where i will eventually write things")
	return s.String()
}

func (m model) resumeView() string {
	if m.width < 83 {
		return m.paraStyle.Render("Please expand terminal out to at least 83 characters")
	}
	return m.viewport.View()
}

func (m model) aboutView() string {
	aboutText := `this ssh application was built using charmbraclet's bubbletea, wish, and lipgloss! you can find the source code for it at cjflan/connorjf.ssh

if you run into any issues while browsing (or weird colorings) feel free to open an issue on the repo and i will do my best to fix it

inspriation for this site came from terminal.shop`

	var s strings.Builder
	s.WriteString(m.paraStyle.SetString(aboutText).String())
	return s.String()
}
