package menubar

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Menubar is a widget that represents a menubar.
type Menubar struct {
	MenuItems []MenuItem
}

// MenuItem is a single item in the menubar.
type MenuItem struct {
	Clickable *widget.Clickable
	OnClick   func(gtx layout.Context)
	Layout    func(gtx layout.Context, th *material.Theme) layout.Dimensions
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
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if item.Clickable != nil && item.Clickable.Clicked(gtx) {
				item.OnClick(gtx)
			}
			if item.Clickable != nil {
				return material.Clickable(gtx, item.Clickable, func(gtx layout.Context) layout.Dimensions {
					return item.Layout(gtx, th)
				})
			}
			return item.Layout(gtx, th)
		})
	}

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx, children...)
}
