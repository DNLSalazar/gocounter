package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	nameInput = iota
	valueInput
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	dartGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(dartGray)
)

func (m *CounterAppModel) setCreateInputs() {
	m.createInfo.inputs[nameInput].SetValue("")
	m.createInfo.inputs[valueInput].SetValue(fmt.Sprint(0))
}

func createInputs() []textinput.Model {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[nameInput] = textinput.New()
	inputs[nameInput].Placeholder = "My counter"
	inputs[nameInput].Focus()
	inputs[nameInput].CharLimit = 30
	inputs[nameInput].Width = 35
	inputs[nameInput].Prompt = ""

	inputs[valueInput] = textinput.New()
	inputs[valueInput].Placeholder = "0"
	inputs[valueInput].CharLimit = 3
	inputs[valueInput].Width = 10
	return inputs
}

func (m *CounterAppModel) createCreatePage() string {
	if m.helpPage {
		return m.createCreateHelpPage()
	} else {
		return m.createCreateContentPage()
	}
}

func (m *CounterAppModel) createCreateContentPage() string {
	s := "Press Ctrl+h for help\r\n"
	s += "Please insert initial name and value\r\n"
	s += "Press Ctrl+c to go back without saving\r\n"
	s += inputStyle.Width(20).Render("Name")
	s += m.createInfo.inputs[nameInput].View()
	s += "\r\n\r\n"
	s += inputStyle.Width(30).Render("Value")
	s += m.createInfo.inputs[valueInput].View()
	return s
}

func (m *CounterAppModel) createCreateHelpPage() string {
	s := "Press Ctrl+h to toggle help\r\n\r\n"
	s += "This page is for creaeting a new counter\r\n"
	s += "Press Tab to move between inputs\r\n"
	s += "Press Enter when you are done to save the new counter\r\n"
	s += "Press Ctrl+c to exit edit mode"
	return s
}

func (m *CounterAppModel) nextCreationInput() {
	if m.createInfo.phase == nameInput {
		m.createInfo.phase = valueInput
	} else {
		m.createInfo.phase = nameInput
	}
}

func (m *CounterAppModel) handleUpdateCreatePage(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					m.dbService.Insert(name, value)
					m.updateData()
					m.selected = &(*m.data)[0]
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
