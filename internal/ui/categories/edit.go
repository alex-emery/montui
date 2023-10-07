package categories

import (
	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/pkg/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

// I'm not a fan of this whole thing, its super gross
// Closing the view on a "GetCategoriesMsg" isn't great

type categoryEditModel struct {
	inputName  textinput.Model
	inputColor textinput.Model
	id         uint
	cursor     int
}

func newCategoryEdit(id uint, category, color string) categoryEditModel {
	inputName := textinput.New()
	inputName.SetValue(category)
	inputName.Focus()
	inputColor := textinput.New()
	inputColor.SetValue(color)

	return categoryEditModel{id: id, inputName: inputName, inputColor: inputColor}
}

func (m categoryEditModel) Init() tea.Cmd { return nil }

func (m categoryEditModel) Update(msg tea.Msg) (categoryEditModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, app.UpdateCategory(storage.Category{Model: gorm.Model{
				ID: m.id,
			}, Name: m.inputName.Value(), Color: m.inputColor.Value()})
		case "up", "down":
			if m.cursor == 0 {
				m.cursor = 1
				m.inputName.Blur()
				return m, m.inputColor.Focus()
			}
			m.cursor = 0
			m.inputColor.Blur()
			return m, m.inputName.Focus()
		}
	}

	m.inputName, cmd = m.inputName.Update(msg)
	cmds = append(cmds, cmd)

	m.inputColor, cmd = m.inputColor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

type CategorySelect struct {
	ID uint
}

func (m categoryEditModel) View() string {
	s := lipgloss.JoinVertical(lipgloss.Left, "Edit Category",
		"Name "+m.inputName.View(),
		"Color "+m.inputColor.View(),
	)

	return s
}
