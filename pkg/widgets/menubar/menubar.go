package menubar

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Menubar is a widget that represents a menubar.
type Menubar struct {
	MenuItems []MenuItem
}

// MenuItem is a single item in the menubar.
type MenuItem interface {
	IsFlexed() bool
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
}

// NewMenubar creates a new menubar.
func NewMenubar() *Menubar {
	return &Menubar{}
}

func (m *Menubar) AddMenuItem(menuItem MenuItem) {
	m.MenuItems = append(m.MenuItems, menuItem)
}

func (menu *Menubar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	children := make([]layout.FlexChild, len(menu.MenuItems))
	for i, item := range menu.MenuItems {
		layoutMenuItem := func(gtx layout.Context) layout.Dimensions {
			return item.Layout(gtx, th)
		}
		if item.IsFlexed() {
			children[i] = layout.Flexed(1, layoutMenuItem)
		} else {
			children[i] = layout.Rigid(layoutMenuItem)
		}
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx, children...)
}
