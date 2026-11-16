package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/lexicon/internal/config"
	"github.com/kristyancarvalho/lexicon/internal/provider"
)

type mode int

const (
	modeChat mode = iota
	modeToken
)

const chromeHeight = 4

type Model struct {
	cfg       *config.Config
	cfgPath   string
	providers []provider.Info
	mode      mode
	width     int
	height    int
	viewport  viewport.Model
	input     textinput.Model
	token     textinput.Model
	history   []provider.Message
	pending   bool
	notice    string
}

func New(cfg *config.Config, cfgPath string) Model {
	input := textinput.New()
	input.Placeholder = "Escreva uma mensagem"
	input.Prompt = "› "
	input.Focus()

	token := textinput.New()
	token.Placeholder = "cole o token aqui"
	token.Prompt = "› "
	token.EchoMode = textinput.EchoPassword
	token.EchoCharacter = '•'
	token.CharLimit = 512

	return Model{
		cfg:       cfg,
		cfgPath:   cfgPath,
		providers: provider.All(),
		viewport:  viewport.New(0, 0),
		input:     input,
		token:     token,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		return m, nil
	case replyMsg:
		m.pending = false
		if msg.err != nil {
			m.notice = msg.err.Error()
		} else {
			m.notice = ""
			m.history = append(m.history, provider.Message{Role: provider.RoleAssistant, Content: msg.content})
		}
		m.refresh()
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.mode == modeToken {
			return m.updateToken(msg)
		}
		return m.updateChat(msg)
	}

	var cmd tea.Cmd
	if m.mode == modeToken {
		m.token, cmd = m.token.Update(msg)
	} else {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m Model) updateChat(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+k":
		m.openToken()
		return m, textinput.Blink
	case "ctrl+p":
		m.cycleProvider()
		return m, nil
	case "ctrl+l":
		m.history = nil
		m.notice = ""
		m.refresh()
		return m, nil
	case "pgup", "pgdown":
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case "enter":
		return m.submit()
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) updateToken(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.closeToken()
		return m, nil
	case "tab":
		m.cycleProvider()
