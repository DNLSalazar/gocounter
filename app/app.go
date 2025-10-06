package app

import (
	"context"
	"time"

	"github.com/DNLSalazar/gocounter/db"
	"github.com/DNLSalazar/gocounter/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	selectPage = iota
	counterPage
	deletePage
	createPage
	editPage
)

type createCounterInfo struct {
	name   string
	value  int
	inputs []textinput.Model
	phase  int
	error  string
}

type CounterAppModel struct {
	data          *[]models.Counter
	page          int
	selected      *models.Counter
	selectedIndex int
	dbService     *db.DatabaseService
	createInfo    createCounterInfo
	helpPage      bool
	ctx           *context.Context
	cancel        *context.CancelFunc
}

func (m *CounterAppModel) saveData() {
	cancel := *(m.cancel)
	cancel()
	m.dbService.SaveFile()
}

func (m *CounterAppModel) counterToSave() {
	ctx := *(m.ctx)
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Second * 5):
		m.dbService.SaveFile()
	}
}

func CreateApp(db *db.DatabaseService) *tea.Program {
	data := db.Get()
	var selected *models.Counter
	if len(data) > 0 {
		selected = &data[0]
	} else {
		selected = nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	model := CounterAppModel{
		data:      &data,
		selected:  selected,
		dbService: db,
		helpPage:  false,
		createInfo: createCounterInfo{
			name:   "",
			value:  0,
			phase:  nameInput,
			inputs: createInputs(),
		},
		ctx:    &ctx,
		cancel: &cancel,
	}

	return tea.NewProgram(model)
}

func (m *CounterAppModel) isDataEmpty() bool {
	return len(*m.data) == 0
}

func (m *CounterAppModel) updateData() {
	newData := m.dbService.Get()
	m.data = &newData
}

func (m CounterAppModel) Init() tea.Cmd {
	return nil
}

func (m CounterAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.page {
	case selectPage:
		return m.handleUpdateSelectPage(msg)
	case counterPage:
		return m.handleUpdateCounterPage(msg)
	case deletePage:
		return m.handleUpdateDeletePage(msg)
	case createPage:
		return m.handleUpdateCreatePage(msg)
	case editPage:
		return m.handleUpdateEditPage(msg)
	default:
		return m, nil
	}
}

func (m CounterAppModel) View() string {
	switch m.page {
	case selectPage:
		return m.createSelectPage()
	case counterPage:
		return m.createCounterPage()
	case deletePage:
		return m.createDeletePage()
	case createPage:
		return m.createCreatePage()
	case editPage:
		return m.createEditPage()
	}
	return "NOT FOUND"
}
