package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


var title_style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#c6a0f6"))

var banner_style = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#c6a0f6"))

var program_style = lipgloss.NewStyle().
    Bold(false).
    Foreground(lipgloss.Color("#a6da95"))

var exit_style = lipgloss.NewStyle().
    Bold(false).
    Foreground(lipgloss.Color("#ed8796"))

var footer_style = lipgloss.NewStyle().
    Bold(false).
    Foreground(lipgloss.Color("#5b6078"))

type model struct {
    run bool
    program string
    width int
    height int
    exit_msg string

}
type FinishedMsg struct { err error }
func initialModel() model {
	return model{
        width: 0,
        height: 0,
    }
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "ctrl+c", "q":
            return m, tea.Quit
        case "enter", " ":
            return m, tea.ExecProcess(exec.Command("zsh", "-ci", m.program), func(err error) tea.Msg {
                return FinishedMsg{err: err}
            })
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case FinishedMsg:
        if msg.err != nil {
            m.exit_msg = "Last run: " + msg.err.Error()
        }
    }

    return m, nil
}

func (m model) View() string {
    // The header
    // text = title_style.Render("Ready when you are...")
    banner_style = banner_style.Width(m.width-2).Height(m.height-2).PaddingTop(m.height/2-1)
    // title_style = title_style.Width(m.width-2)
    // program_style = program_style.Width(m.width-2)
    // footer_style = footer_style.Width(m.width-2)

    center := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width-2)

    text := lipgloss.JoinVertical(lipgloss.Center, 
        lipgloss.JoinHorizontal(lipgloss.Center, 
            title_style.Render("Waiting to run: "),
            program_style.Render(m.program),
        ),
        exit_style.Render("\n" + m.exit_msg),
        footer_style.Render("\nPress [Enter] to continue... [Ctrl+C | q] to cancel"),
    )

    text = center.Render(text)
    return banner_style.Render(text)
    
    
}

func main() {
    program := strings.Join(os.Args[1:], " ")

    m := model{
        run: false,
        program: program,
        width: 0,
        height: 0,
    }
    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
