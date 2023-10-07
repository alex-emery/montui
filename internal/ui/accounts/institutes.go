package accounts

import (
	"fmt"
	"io"
	"strings"

	"github.com/alex-emery/montui/internal/ui/app"
	"github.com/alex-emery/montui/pkg/nordigen"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type InstituteModel struct {
	list list.Model
}

type item struct {
	nordigen.Integration
}

func (i item) FilterValue() string { return i.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func newInstitutionPicker() *InstituteModel {
	return &InstituteModel{list: list.New(nil, itemDelegate{}, 0, 0)}
}

func (m InstituteModel) Init() tea.Cmd {
	return app.GetBanks()
}

func bankToItems(banks []nordigen.Integration) []list.Item {
	var items = make([]list.Item, 0, len(banks))
	for _, bank := range banks {
		items = append(items, item{
			bank,
		})
	}

	return items
}

func (m InstituteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case app.NewBanksMsg:
		items := bankToItems(msg.Banks)
		cmd := m.list.SetItems(items)
		return m, cmd
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-20, msg.Height-20)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.list.SelectedItem().(item)
			return m, app.SelectBank(selected.Integration)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m InstituteModel) View() string {
	return m.list.View()
}
