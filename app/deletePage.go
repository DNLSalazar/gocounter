package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *CounterAppModel) createDeletePage() string {
	s := ""
	s += fmt.Sprintf("Are you sure you want to delete counter %v?\r\n", m.selected.Name)
	s += fmt.Sprintf("Press Y[Yes] / N[No]")
	return s
}

func (m *CounterAppModel) handleUpdateDeletePage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.helpPage = false
			return m, tea.Quit
		case tea.KeyEsc:
			m.page = selectPage
		}
		switch msg.String() {
		case "y", "Y":
			m.dbService.Delete(m.selected.Id)
			m.updateData()
			m.page = selectPage
			m.selectedIndex = 0
			if m.isDataEmpty() {
				m.selected = nil
			} else {
				m.selected = &(*m.data)[m.selectedIndex]
			}
		case "n", "N":
			m.page = selectPage
		}
	}
	m.helpPage = false
	return *m, nil
}
