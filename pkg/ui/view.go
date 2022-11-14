package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/koki-develop/clive/pkg/util"
)

func (m *Model) View() string {
	if m.err != nil {
		return m.errView()
	}

	if m.config == nil {
		return m.loadingConfigView()
	}

	if m.page == nil {
		return m.openingView()
	}

	s := m.actionsView()
	if m.quitting {
		s += "\n\n" + m.quittingView()
	}

	return s
}

func (m *Model) errView() string {
	if !m.running() {
		return ""
	}

	return styleErrorHeader.Render("Error") + "\n" + m.err.Error() + "\n\n" + m.quittingView()
}

func (m *Model) loadingConfigView() string {
	return fmt.Sprintf("%s Loading config", m.spinner.View())
}

func (m *Model) openingView() string {
	return fmt.Sprintf("%s Opening", m.spinner.View())
}

func (m *Model) actionsView() string {
	from := util.Max(0, m.currentActionIndex-3)
	show := 20
	digits := util.Digits(len(m.config.Actions))

	rows := []string{}
	for i, action := range m.config.Actions {
		if i < from && len(m.config.Actions)-i > show {
			continue
		}
		if i-from >= show {
			rows = append(rows, fmt.Sprintf("... %d more actions", len(m.config.Actions)-i))
			break
		}

		var style lipgloss.Style

		cursor := "  "
		if m.currentActionIndex > i {
			style = styleDone
		} else if m.currentActionIndex == i {
			style = styleActive
			if !m.quitting {
				if m.pausing {
					cursor = "> "
				} else {
					cursor = m.spinner.View()
				}
			}
		}

		s, trunc := util.TruncateString(action.String(), 40)
		if trunc {
			s += styleTruncated.Render("...")
		}

		num := util.PaddingRight(fmt.Sprintf("#%d", i+1), digits+1)
		rows = append(rows, fmt.Sprintf("%s %s%s", style.Render(num), cursor, style.Render(s)))
	}

	return styleActionHeader.Render("Actions") + "\n" + strings.Join(rows, "\n")
}

func (m *Model) quittingView() string {
	return styleActive.Render("Press enter to quit")
}
