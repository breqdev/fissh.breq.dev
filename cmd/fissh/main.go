package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/breqdev/fissh.breq.dev/internal/fishes"
	"github.com/breqdev/fissh.breq.dev/internal/timezone"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
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
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	renderer := bubbletea.MakeRenderer(s)
	appStyle := renderer.NewStyle()

	timezone := timezone.LookupTimezone(s.RemoteAddr().String())
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(err)
	}

	m := model{
		timer:    timer.NewWithInterval(999999999*time.Second, time.Millisecond),
		time:     time.Now(),
		timezone: loc,
		page:     HomePage,
		fish:     "",
		window:   tea.WindowSizeMsg{Width: pty.Window.Width, Height: pty.Window.Height},
		styles: appStyles{
			app:  appStyle,
			txt:  renderer.NewStyle().Foreground(lipgloss.Color("2")).Inherit(appStyle),
			fish: renderer.NewStyle().Foreground(lipgloss.Color("4")).Inherit(appStyle),
			quit: renderer.NewStyle().Foreground(lipgloss.Color("8")).Inherit(appStyle),
		},
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type appStyles struct {
	app  lipgloss.Style
	txt  lipgloss.Style
	fish lipgloss.Style
	quit lipgloss.Style
}

const (
	HomePage = iota
	AboutPage
)

type model struct {
	timer    timer.Model
	time     time.Time
	timezone *time.Location
	page     int
	fish     string
	window   tea.WindowSizeMsg
	styles   appStyles
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.page = AboutPage
			return m, nil
		case "esc":
			m.page = HomePage
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.window = msg
		return m, nil
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.time = time.Now()
		if m.time.In(m.timezone).Format("03:04") == "11:11" || os.Getenv("ALWAYSFISH") == "1" {
			if m.fish == "" {
				m.fish = fishes.GetFish(m.window.Width, m.window.Height)
			}
		}
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.page {
	case HomePage:
		return m.HomePage()
	case AboutPage:
		return m.AboutPage()
	default:
		return m.styles.txt.Render("Page Not Found")
	}
}

func (m model) HomePage() string {
	s := fmt.Sprintf("the time is %s", m.time.In(m.timezone).Format("15:04:05"))
	if m.fish != "" {
		s = fmt.Sprintf("%s\n\n%s\n\n%s", m.styles.txt.Render(s), m.styles.fish.Render(m.fish), m.styles.fish.Render("make a fish"))
	} else {
		s = fmt.Sprintf("%s\n\n%s", m.styles.txt.Render(s), m.styles.fish.Render("come back at 11:11"))
	}
	s = fmt.Sprintf("%s\n\n%s", s, m.styles.quit.Render("press 'a' for about or 'q' to quit"))

	s = lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center, s)

	return m.styles.app.Width(m.window.Width).Height(m.window.Height).Render(s)
}

func (m model) AboutPage() string {
	aboutContent := fmt.Sprintf("%s\n\n%s", m.styles.txt.Render("Made with <3 by @breqdev and @avasilver"), m.styles.quit.Render("Press 'esc' to go back or 'q' to quit"))
	return lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center, aboutContent)
}
