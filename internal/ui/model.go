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

const chromeHeight = 6

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
		return m, nil
	case "enter":
		return m.saveToken()
	}

	var cmd tea.Cmd
	m.token, cmd = m.token.Update(msg)
	return m, cmd
}

func (m Model) saveToken() (tea.Model, tea.Cmd) {
	value := strings.TrimSpace(m.token.Value())
	if value == "" {
		m.notice = "token vazio"
		return m, nil
	}
	m.cfg.SetToken(m.active().ID, value)
	if err := m.cfg.Save(m.cfgPath); err != nil {
		m.notice = err.Error()
		return m, nil
	}
	m.closeToken()
	return m, nil
}

func (m Model) submit() (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(m.input.Value())
	if text == "" || m.pending {
		return m, nil
	}

	info := m.active()
	p, err := provider.New(info.ID, m.cfg.Token(info.ID))
	if err != nil {
		m.notice = err.Error()
		return m, nil
	}

	m.history = append(m.history, provider.Message{Role: provider.RoleUser, Content: text})
	m.input.Reset()
	m.notice = ""
	m.pending = true
	m.refresh()
	return m, sendCmd(p, m.cfg.Model(info.ID, info.DefaultModel), m.history)
}

func (m Model) active() provider.Info {
	if info, ok := provider.Lookup(m.cfg.Active); ok {
		return info
	}
	return m.providers[0]
}

func (m *Model) cycleProvider() {
	current := m.active().ID
	next := m.providers[0]
	for i, info := range m.providers {
		if info.ID == current {
			next = m.providers[(i+1)%len(m.providers)]
			break
		}
	}
	m.cfg.Active = next.ID
	if err := m.cfg.Save(m.cfgPath); err != nil {
		m.notice = err.Error()
	}
}

func (m *Model) openToken() {
	m.mode = modeToken
	m.notice = ""
	m.token.Reset()
	m.input.Blur()
	m.token.Focus()
}

func (m *Model) closeToken() {
	m.mode = modeChat
	m.notice = ""
	m.token.Reset()
	m.token.Blur()
	m.input.Focus()
}

func (m *Model) layout() {
	m.viewport.Width = m.width
	m.viewport.Height = max(m.height-chromeHeight, 1)
	m.input.Width = max(m.width-4, 10)
	m.token.Width = 40
	m.refresh()
}

func (m *Model) refresh() {
	m.viewport.SetContent(m.renderHistory())
	m.viewport.GotoBottom()
}
