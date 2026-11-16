package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/lexicon/internal/provider"
)

const requestTimeout = 2 * time.Minute

type replyMsg struct {
	content string
	err     error
}

func sendCmd(p provider.Provider, model string, history []provider.Message) tea.Cmd {
	messages := append([]provider.Message(nil), history...)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		content, err := p.Chat(ctx, model, messages)
		return replyMsg{content: content, err: err}
	}
}
