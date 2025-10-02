package main

import (
	"fmt"
	"log"
	// "os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	// payload string
	ti textinput.Model // text input model will have its own view, method, and etc methods
}

func initialModel() model {
	ti := textinput.New()
	return model{
		ti,
	}

}

// automatically run when p.Run()
func (m model) Init() tea.Cmd {
	return nil
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

func main() {
	//fmt.Println("foo bar")
	p := tea.NewProgram(initialModel()) // create new program, tea.NewProgram() returns a pointer to a pointer to a t.program structure

	// p.Run() returns: t.model and error
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
