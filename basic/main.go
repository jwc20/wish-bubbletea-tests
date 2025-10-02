package main 


import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	fmt.Println("foo bar")
	p := tea.NewProgram() // create new program, tea.NewProgram() returns a pointer to a pointer to a t.program structure
	

	// p.Run() returns: t.model and error
	if _, err := p. Run(); err != nil {
		log.Fatal(err)
	}
}
