package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = "3000"
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

	// go routine to handle ssh server
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

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	s.Pty()
	return initialModel(), []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	// payload string
	ti textinput.Model // text input model will have its own view, method, and etc methods
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Jae C"
	ti.Width = 20
	return model{
		ti,
	}

}

/* --------------------------------------------------------- */

// automatically run when p.Run()
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// automatically run when p.Run()
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// this method returns tea.Model beacuause this is not a pointer/receiver
	// any changes made to m model will not persist outside of this method scope because it's passed by copy
	// this meathod is like an event handler (pub/sub ood pattern) where it listens for events (in the form of t.message)
	// return m, nil

	if val, ok := msg.(tea.KeyMsg); ok {
		key := val.String()
		// os.WriteFile("output.log", []byte(key), 0644)

		if key == "ctrl+c" {
			return m, tea.Quit
		}
		if key == "enter" {
			// save to file
			os.WriteFile("output.log", []byte(m.ti.Value()), 0644)
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)

	return m, cmd
}

func (m model) View() string {
	// return m.payload
	// return m.ti.View()
	output := fmt.Sprintf("Name?\n\n%v", m.ti.View())
	return output
}

//func main() {
//	//fmt.Println("foo bar")
//	p := tea.NewProgram(initialModel(), tea.WithAltScreen()) // create new program, tea.NewProgram() returns a pointer to a pointer to a t.program structure
//
//	// p.Run() returns: t.model and error
//	if _, err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
//}
