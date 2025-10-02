package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

// This example was inspired by terminal.shop - a TUI-based coffee shop accessible over SSH
// The goal is to create a simple "Hello World" app using the same technologies

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
	// For production deployment, use 0.0.0.0 to listen on all interfaces
	// localhost is good for development
	host = "0.0.0.0"
	// Port 22 is the default SSH port but requires elevated privileges
	// Using port 3000 instead to avoid permission issues on macOS
	port = "3000"
)

func main() {
	// Wish handles all SSH security, user management, and shell restrictions
	// This prevents users from gaining shell or root access to the server
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		// SSH keys will be stored in .ssh/id_ed25519
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			// The bubbletea middleware connects our TUI app to SSH sessions
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	// Go routine (similar to multi-threading) to handle ssh server in parallel
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

// teaHandler is called for each SSH connection
// In a Wish app, you don't call tea.NewProgram().Run() directly
// Instead, you return the model and options to the middleware
// The middleware handles running, stopping, and managing the program
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// PTY (pseudo-terminal) can provide info about client's terminal
	// (terminal width, height, color scheme, etc.) but we're not using it here
	s.Pty()
	// WithAltScreen makes the app take over the entire terminal screen
	// Similar to how terminal.shop creates a full-screen experience
	return initialModel(), []tea.ProgramOption{tea.WithAltScreen()}
}

// Model represents the state of the entire app (following Elm architecture)
// Bubble Tea is immutable - we update by returning a new model with changes
type model struct {
	// payload string
	// Using a pre-built text input component from Bubbles (component library)
	// The text input has its own update, view, and init methods
	ti textinput.Model // text input model will have its own view, method, and etc methods
}

// Constructor for creating the initial model state
func initialModel() model {
	ti := textinput.New()
	// Focus is important - without it, the text input won't respond to typing
	// Multiple text inputs can exist, but only the focused one receives input
	ti.Focus()
	ti.Placeholder = "Jae C"
	// Width must be set for placeholder to display correctly
	ti.Width = 20
	return model{
		ti,
	}

}

/* --------------------------------------------------------- */

// Init is automatically called by Bubble Tea when the program starts
// We never call this directly - Bubble Tea calls it for us
func (m model) Init() tea.Cmd {
	// Blink command makes the cursor start blinking immediately
	// Without this, cursor would be static until first keystroke
	return textinput.Blink
}

// Update is the event handler - called automatically when messages (events) occur
// This is not a pointer receiver, so changes aren't persisted unless returned
// Similar to React's immutable state updates
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// this method returns tea.Model beacuause this is not a pointer/receiver
	// any changes made to m model will not persist outside of this method scope because it's passed by copy
	// this meathod is like an event handler (pub/sub ood pattern) where it listens for events (in the form of t.message)
	// return m, nil

	// Type assertion to check if the message is a keyboard event
	if val, ok := msg.(tea.KeyMsg); ok {
		// String() method returns string representation of the key pressed
		key := val.String()
		// os.WriteFile("output.log", []byte(key), 0644)

		// Without handling ctrl+c, the app becomes unresponsive
		// Users would need to kill the process manually (e.g., using htop)
		if key == "ctrl+c" {
			// tea.Quit tells Bubble Tea to stop the application
			return m, tea.Quit
		}
		if key == "enter" {
			// save to file
			// ti.Value() gets the current text from the input field
			// 0644 is octal file permission: read/write for owner, read for group/others
			os.WriteFile("output.log", []byte(m.ti.Value()), 0644)
			return m, tea.Quit
		}
	}

	// Pass the message to the text input component for processing
	// The text input returns its updated model and any commands
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)

	// Return the updated model with the new text input state
	// Commands from text input are forwarded to Bubble Tea
	return m, cmd
}

// View renders the UI - returns a string that appears in the terminal
// Called automatically whenever the model changes
func (m model) View() string {
	// return m.payload
	// return m.ti.View()
	// fmt.Sprintf creates a formatted string with the prompt and input field
	output := fmt.Sprintf("Name?\n\n%v", m.ti.View())
	return output
}
