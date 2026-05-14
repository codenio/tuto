package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Border(lipgloss.NormalBorder()).
			Padding(0, 1).
			BorderForeground(lipgloss.Color("62"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			BorderForeground(lipgloss.Color("241")).
			Width(72)

	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	cmdStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51"))
)

// Title renders a heading.
func Title(s string) string {
	return titleStyle.Render(s)
}

// Box wraps body text in a bordered block.
func Box(body string) string {
	return boxStyle.Render(strings.TrimSpace(body))
}

// Success message.
func Success(s string) string {
	return okStyle.Render(s)
}

// Error message.
func Error(s string) string {
	return errStyle.Render(s)
}

// Muted line.
func Muted(s string) string {
	return mutedStyle.Render(s)
}

// RenderMarkdown renders s as markdown for the terminal.
// Falls back to plain text if stdout is not a TTY or if glamour fails.
// WithAutoStyle is safe here because we already guard on IsTerminal — it
// will never fall back to notty when stdout is a real TTY.
func RenderMarkdown(s string) string {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return s
	}
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(64),
	)
	if err != nil {
		return s
	}
	rendered, err := renderer.Render(s)
	if err != nil {
		return s
	}
	return strings.TrimSpace(rendered)
}

// FormatStepBody builds a plain-text step body (no glamour) for use inside a lipgloss Box.
// Kept for callers that need an inline string. Use PrintStep for the full rendered display.
func FormatStepBody(instruction, commandToRun string) string {
	inst := strings.TrimSpace(instruction)
	cmd := strings.TrimSpace(commandToRun)
	var b strings.Builder
	if inst != "" {
		b.WriteString(inst)
	}
	if cmd == "" {
		return b.String()
	}
	if inst != "" {
		b.WriteString("\n\n")
	}
	b.WriteString(mutedStyle.Render("Command `tuto step next` will run:"))
	b.WriteString("\n")
	b.WriteString(cmdStyle.Render(cmd))
	return b.String()
}

// PrintStep renders a full step block: header line, glamour-rendered instruction, command.
func PrintStep(stepNum, total int, id, instruction, commandToRun string) {
	header := fmt.Sprintf("Step %d / %d  ·  %s", stepNum, total, id)
	fmt.Println(mutedStyle.Render("─── " + header + " "))
	if inst := strings.TrimSpace(instruction); inst != "" {
		fmt.Println(RenderMarkdown(inst))
	}
	if cmd := strings.TrimSpace(commandToRun); cmd != "" {
		fmt.Println(mutedStyle.Render("  Validate with:  ") + cmdStyle.Render(cmd))
	}
	fmt.Println()
}

// ProgressBar draws a simple ASCII progress bar [=====     ].
func ProgressBar(done, total int, width int) string {
	if total <= 0 {
		return "[----------]"
	}
	if done < 0 {
		done = 0
	}
	if done > total {
		done = total
	}
	if width < 8 {
		width = 20
	}
	inner := width - 2
	filled := (done * inner) / total
	if done == total && total > 0 {
		filled = inner
	}
	sb := strings.Builder{}
	sb.WriteByte('[')
	for i := 0; i < inner; i++ {
		if i < filled {
			sb.WriteByte('=')
		} else {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// ProgressLine returns a labeled progress line.
func ProgressLine(label string, done, total int) string {
	pct := 0
	if total > 0 {
		pct = (done * 100) / total
	}
	return fmt.Sprintf("%s %s %d%%", label, ProgressBar(done, total, 22), pct)
}
