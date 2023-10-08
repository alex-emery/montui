package rules

import (
	"fmt"

	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/pkg/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type editModel struct {
	inputPattern  textinput.Model
	inputCategory textinput.Model
	id            uint
	cursor        int
	edit          bool
}

func newEditModel(id uint, pattern, category string, edit bool) editModel {
	inputPattern := textinput.New()
	inputPattern.SetValue(pattern)
	inputPattern.Focus()
	inputCategory := textinput.New()
	inputCategory.SetValue(category)

	return editModel{id: id, inputPattern: inputPattern, inputCategory: inputCategory, edit: edit}
}

func (m editModel) Init() tea.Cmd { return nil }

func (m editModel) Update(msg tea.Msg) (editModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(m.inputPattern.Value()) == 0 {
				return m, app.SendError(fmt.Errorf("pattern can not be empty"))
			}
			if len(m.inputCategory.Value()) == 0 {
				return m, app.SendError(fmt.Errorf("category can not be empty"))
			}

			rule := &storage.Rule{
				Model: gorm.Model{
					ID: m.id,
				},
				Pattern: m.inputPattern.Value(),
				Category: storage.Category{
					Name: m.inputCategory.Value(),
				}}
			if m.edit {
				return m, app.UpdateRule(rule)
			}
			return m, app.CreateRule(rule)

		case "up", "down":
			if m.cursor == 0 {
				m.cursor = 1
				m.inputPattern.Blur()
				return m, m.inputCategory.Focus()
			}
			m.cursor = 0
			m.inputCategory.Blur()
			return m, m.inputPattern.Focus()
		}
	}

	m.inputPattern, cmd = m.inputPattern.Update(msg)
	cmds = append(cmds, cmd)

	m.inputCategory, cmd = m.inputCategory.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m editModel) View() string {
	s := lipgloss.JoinVertical(lipgloss.Left, "Edit Category",
		"Pattern "+m.inputPattern.View(),
		"Category "+m.inputCategory.View(),
	)

	return s
}
