package app

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *CounterAppModel) setEditInputs() {
	m.createInfo.inputs[nameInput].SetValue(m.selected.Name)
	m.createInfo.inputs[valueInput].SetValue(fmt.Sprint(m.selected.Value))
}

func (m *CounterAppModel) createEditPage() (s string) {
	if m.helpPage {
		return m.createEditHelpPage()
	} else {
		return m.createEditContentPage()
	}
}

func (m *CounterAppModel) createEditContentPage() string {
	s := "Press Ctrl+h for help\r\n"
	s += "Please insert new name and value\r\n"
	s += "Press Ctrl+c to go back without saving\r\n"
	s += inputStyle.Width(20).Render("Name")
	s += m.createInfo.inputs[nameInput].View()
	s += fmt.Sprint("\r\n\r\n")
	s += inputStyle.Width(30).Render("Value")
	s += m.createInfo.inputs[valueInput].View()
	return s
}

func (m *CounterAppModel) createEditHelpPage() string {
	s := "Press Ctrl+h to toggle help\r\n\r\n"
	s += "This page is for editing the counter\r\n"
	s += "Press Tab to move between inputs\r\n"
	s += "Press Enter when you are done editing the counter\r\n"
	s += "Press Ctrl+c to exit edit mode"
	return s
}

func (m *CounterAppModel) handleUpdateEditPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC:
			m.helpPage = false
			m.page = selectPage
			return *m, nil
		case tea.KeyCtrlH:
			m.helpPage = !m.helpPage
			return *m, nil
		case tea.KeyTab:
			m.nextCreationInput()

		case tea.KeyEnter:
			if m.createInfo.phase == 1 {
				name := strings.TrimSpace(m.createInfo.inputs[nameInput].Value())
				value, err := strconv.Atoi(m.createInfo.inputs[valueInput].Value())
				if err != nil {
					m.createInfo.error = "Error on value, please introduce a valid number"
					return *m, nil
				}
				if name != "" {
					m.createInfo.error = ""
					m.dbService.Update(m.selected.Id, name, value)
					m.updateData()
				}
				m.createInfo.phase = nameInput
				m.page = selectPage
				return *m, nil
			}
			m.nextCreationInput()
		}
		var cmds []tea.Cmd = make([]tea.Cmd, len(m.createInfo.inputs))
		for i := range (*m).createInfo.inputs {
			m.createInfo.inputs[i], cmds[i] = m.createInfo.inputs[i].Update(msg)
			m.createInfo.inputs[i].Blur()
		}
		m.createInfo.inputs[m.createInfo.phase].Focus()
		return *m, tea.Batch(cmds...)
	}
	return *m, nil
}
