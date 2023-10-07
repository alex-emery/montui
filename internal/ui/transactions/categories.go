package transactions

import (
	"fmt"

	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type categoryPicker struct {
	height     int
	width      int
	categories []storage.Category
	cursor     int
}

func (m categoryPicker) Init() tea.Cmd { return nil }

func (m categoryPicker) Update(msg tea.Msg) (categoryPicker, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case app.NewCategoriesMsg:
		m.categories = msg.Categories
	case tea.WindowSizeMsg:
		m.height = msg.Height - 10
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.categories)-1 {
				m.cursor++
			}
		case "enter":
			return m, SelectCategory(m.categories[m.cursor].ID)
		}

	}

	return m, cmd
}
func SelectCategory(categoryID uint) tea.Cmd {
	return func() tea.Msg {
		return CategorySelect{
			ID: categoryID,
		}
	}
}

type CategorySelect struct {
	ID uint
}

func (m categoryPicker) View() string {
	s := "Categories\n"
	for i, choice := range m.categories {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if i == m.cursor {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Name)
	}

	return s
}

func (m categoryPicker) Cursor() int {
	return m.cursor
}

func (m *categoryPicker) SetCursor(cursor int) {
	m.cursor = cursor
}
