package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *CounterAppModel) createCounterPage() string {
	if m.helpPage {
		return m.createCounterHelpPage()
	} else {
		return m.createCounterContentPage()
	}
}

func (m *CounterAppModel) createCounterContentPage() string {
	s := "Press Ctrl+h for help\r\n"
	s += fmt.Sprintf("%v: %v", m.selected.Name, m.selected.Value)
	return s
}

func (m *CounterAppModel) createCounterHelpPage() string {
	s := "Press Ctrl+h to toggle help\r\n\r\n"
	s += "This page is for incrementing or decrementing a counter value\r\n"
	s += "To increment the counter press up, j or Ctrl+a\r\n"
	s += "To Decrement the counter press down or k\r\n"
	s += "Press Ctrl+c, Enter or q to exit counter mode"
	return s
}

func (m *CounterAppModel) handleUpdateCounterPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.helpPage = false
			m.page = selectPage
		case tea.KeyCtrlC:
			m.helpPage = false
			return m, tea.Quit
		case tea.KeyCtrlH:
			m.helpPage = !m.helpPage
			return *m, nil
		case tea.KeyUp:
			m.dbService.Update(m.selected.Id, m.selected.Name, m.selected.Value+1)
			m.updateData()
			go m.counterToSave()
		case tea.KeyCtrlA:
			m.dbService.Update(m.selected.Id, m.selected.Name, m.selected.Value+1)
			m.updateData()
			go m.counterToSave()
		case tea.KeyDown:
			m.dbService.Update(m.selected.Id, m.selected.Name, m.selected.Value-1)
			m.updateData()
			go m.counterToSave()
		}
		switch msg.String() {
		case "q", "Q":
			m.helpPage = false
			m.page = selectPage
		case "k", "K":
			m.dbService.Update(m.selected.Id, m.selected.Name, m.selected.Value-1)
			m.updateData()
			go m.counterToSave()
		case "j", "J":
			m.dbService.Update(m.selected.Id, m.selected.Name, m.selected.Value+1)
			m.updateData()
			go m.counterToSave()
		}
	}
	return *m, nil
}
