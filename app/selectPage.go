package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const listSpace = "    "
const selected = " -> "

func (m *CounterAppModel) createSelectPage() string {
	if m.helpPage {
		return m.createSelectHelpPage()
	} else {
		return m.createSelectContentPage()
	}
}

func (m *CounterAppModel) createSelectContentPage() string {
	s := fmt.Sprint("Press Ctrl+h for toggle help or q to go back/quit\r\n")
	for _, v := range *m.data {
		space := listSpace
		if m.selected.Id == v.Id {
			space = selected
		}
		s += fmt.Sprintf("%v%v: %v\r\n", space, v.Name, v.Value)
	}
	if len(*m.data) == 0 {
		s += "No entries found! Press A to start creating counters"
	}

	return s
}

func (m *CounterAppModel) createSelectHelpPage() string {
	s := "Press Ctrl+h to toggle help\r\n\r\n"
	s += "This page is for manage your counter\r\n"
	s += "Press a to create a new counter\r\n"
	s += "Press e to edit an existing counter\r\n"
	s += "Press d to delete an existing counter\r\n"
	s += "Press q or Ctrl+c to exit"
	return s
}

func (m *CounterAppModel) updateSelected(n int) {
	if m.isDataEmpty() {
		return
	}
	newIndex := m.selectedIndex + n
	if newIndex < 0 {
		newIndex = len(*m.data) - 1
	} else if newIndex == len(*m.data) {
		newIndex = 0
	}

	m.selectedIndex = newIndex
	m.selected = &(*m.data)[newIndex]
}

func (m *CounterAppModel) handleUpdateSelectPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.isDataEmpty() {
				return *m, nil
			}
			m.page = counterPage
		case tea.KeyUp:
			m.updateSelected(-1)
		case tea.KeyDown:
			m.updateSelected(1)
		case tea.KeyCtrlC:
			m.dbService.SaveFile()
			return m, tea.Quit
		case tea.KeyCtrlH:
			m.helpPage = !m.helpPage
		}
		switch msg.String() {
		case "q", "Q":
			m.dbService.SaveFile()
			tea.Quit()
			return *m, tea.Quit
		case "a", "A":
			m.page = createPage
			m.setCreateInputs()
			return *m, nil
		case "e", "E":
			if m.isDataEmpty() {
				return *m, nil
			}
			m.page = editPage
			m.setEditInputs()
			return *m, nil
		case "d", "D":
			if m.isDataEmpty() {
				return *m, nil
			}
			m.page = deletePage
			return *m, nil
		case "j", "J":
			m.updateSelected(1)
		case "k", "K":
			m.updateSelected(-1)
		}
	}
	return *m, nil
}
