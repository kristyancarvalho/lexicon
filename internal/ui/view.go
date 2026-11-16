package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kristyancarvalho/lexicon/internal/provider"
)

type hint struct{ key, label string }

var (
	chatHints = []hint{
		{"enter", "enviar"}, {"ctrl+k", "token"}, {"ctrl+p", "provider"},
		{"ctrl+l", "limpar"}, {"pgup/pgdn", "rolar"}, {"ctrl+c", "sair"},
	}
	tokenHints = []hint{{"enter", "salvar"}, {"tab", "trocar provider"}, {"esc", "cancelar"}}
)

const welcome = "Bem-vindo ao lexicon. Pressione ctrl+k para configurar o token do provider."

func renderHints(hints []hint) string {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = keyStyle.Render(h.key) + " " + mutedStyle.Render(h.label)
	}
	return strings.Join(parts, mutedStyle.Render("  ·  "))
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}
	if m.mode == modeToken {
		return m.tokenView()
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.header(),
		m.viewport.View(),
		m.rule(),
		m.input.View(),
		m.statusLine(),
		renderHints(chatHints),
	)
}

func (m Model) rule() string {
	return ruleStyle.Render(strings.Repeat("─", m.width))
}

func (m Model) header() string {
	info := m.active()
	model := m.cfg.Model(info.ID, info.DefaultModel)
	dot := okStyle.Render("●")
	if m.cfg.Token(info.ID) == "" {
		dot = errorStyle.Render("●")
	}
	left := titleStyle.Render("lexicon") + "  " + dot + " " + accentStyle.Render(info.Name) + mutedStyle.Render(" · "+model)
	return left + "\n" + m.rule()
}

func (m Model) statusLine() string {
	switch {
	case m.notice != "":
		return errorStyle.Render("✗ " + m.notice)
	case m.pending:
		return mutedStyle.Render("… aguardando resposta")
	}
	return " "
}

func (m Model) renderHistory() string {
	if len(m.history) == 0 {
		return mutedStyle.Render(welcome)
	}

	width := max(m.width-4, 10)
	var b strings.Builder
	for _, msg := range m.history {
		label, bar := userStyle.Render("você"), userBarStyle
		if msg.Role == provider.RoleAssistant {
			label, bar = assistantStyle.Render("lexicon"), botBarStyle
		}
		content := lipgloss.NewStyle().Width(width).Render(msg.Content)
		b.WriteString(bar.Render(label + "\n" + content))
		b.WriteString("\n\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m Model) tokenView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Token de API"))
	b.WriteString("\n\n")

	active := m.active().ID
	for _, info := range m.providers {
		mark := mutedStyle.Render("✗")
		if m.cfg.Token(info.ID) != "" {
			mark = okStyle.Render("✓")
		}
		name := mutedStyle.Render("  " + info.Name)
		if info.ID == active {
			name = activeStyle.Render("▸ " + info.Name)
		}
		b.WriteString(name + " " + mark + "\n")
	}

	b.WriteString("\n")
	b.WriteString(m.token.View())
	b.WriteString("\n\n")
	if m.notice != "" {
		b.WriteString(errorStyle.Render(m.notice) + "\n\n")
	}
	b.WriteString(renderHints(tokenHints))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modalStyle.Render(b.String()))
}
