package cli

import (
	"fmt"
	"schat"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

type chatModel struct {
	Name        string
	viewport    viewport.Model
	messages    []string
	reciveMsg   chan schat.MsgForm
	sendMsg     chan schat.MsgForm
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func NewChatModel(reciveMsg, sendMsg chan schat.MsgForm) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chatModel{
		textarea:    ta,
		messages:    []string{},
		reciveMsg:   reciveMsg,
		sendMsg:     sendMsg,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			msg := schat.MsgForm{
				Name: m.Name,
				Msg:  m.textarea.Value(),
			}
			m.sendMessage(msg)
		}

	case NewMsg:
		message := <-m.reciveMsg
		m.addMsg(message)
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m chatModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m *chatModel) addMsg(message schat.MsgForm) {
	m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", message.Name))+message.Msg)
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
}

func (m *chatModel) sendMessage(message schat.MsgForm) {
	m.sendMsg <- message
}
