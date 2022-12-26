package cli

import (
	"schat"

	tea "github.com/charmbracelet/bubbletea"
)

type status int
type NewMsg struct{}

const (
	logIn status = iota
	chat
	list
)

type mainModel struct {
	focused     status
	readyToChat chan struct{}
	logInModel
	chatModel
}

func NewMainModel(readyToChat chan struct{}, reciveMsg, sendMsg chan schat.MsgForm) *mainModel {
	// LogInView := tea.NewProgram(NewLogInModel())
	// if _, err := LogInView.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	// chatView := tea.NewProgram(cli.NewChatModel())
	// if _, err := chatView.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	return &mainModel{
		focused:     0,
		readyToChat: readyToChat,
		logInModel:  NewLogInModel(),
		chatModel:   NewChatModel(reciveMsg, sendMsg),
	}
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

Rewrite:
	switch m.focused {
	case logIn:
		// switch state
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				name := m.logInModel.textInput.Value()
				m.chatModel.Name = name
				m.focused = chat
				m.readyToChat <- struct{}{}
				break Rewrite
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}
		}
		// handle msg
		var ok bool
		newLogInModel, newCmd := m.logInModel.Update(msg)
		m.logInModel, ok = newLogInModel.(logInModel)
		if !ok {
			panic("failed assertion logInModel")
		}
		cmd = newCmd
	case chat:
		// switch state
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}
			// case NewMsg:
			// 	m.chatModel.
		}
		// handle msg
		var ok bool
		newChatModel, newCmd := m.chatModel.Update(msg)
		m.chatModel, ok = newChatModel.(chatModel)
		if !ok {
			panic("failed assertion chatModel")
		}
		cmd = newCmd
	}
	return m, cmd
}

func (m *mainModel) View() string {
	switch m.focused {
	case logIn:
		return m.logInModel.View()
	case chat:
		return m.chatModel.View()
	}
	return "view err, no view - for you, pls quit"
}
