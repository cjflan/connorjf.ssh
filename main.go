package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	_ "embed"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

//go:embed resume.md
var r string

const (
	host = "localhost"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	renderer := bubbletea.MakeRenderer(s)
	headerStyle := renderer.NewStyle().Foreground(lipgloss.Color("#bb9af7"))
	navStyle := renderer.NewStyle().Foreground(lipgloss.Color("#c0caf5"))
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("#565f89"))
	paraStyle := renderer.NewStyle().Foreground(lipgloss.Color("#c0caf5")).
		Width(80)
	titleStyle := renderer.NewStyle().Foreground(lipgloss.Color("#ff9e64")).
		Bold(true)

	const width = 80
	const height = 20
	vp := viewport.New(width, height)
	vp.Style = renderer.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	vpRenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(TokyoNightStyleConfig),
		glamour.WithWordWrap(width),
	)

	str, _ := vpRenderer.Render(r)

	vp.SetContent(str)

	m := model{
		profile:   renderer.ColorProfile().Name(),
		headerStyle: headerStyle,
		navStyle: navStyle,
		titleStyle: titleStyle,
		txtStyle:  txtStyle,
		quitStyle: quitStyle,
		paraStyle: paraStyle,
		viewport: vp,
		page: home,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
