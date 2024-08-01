package menubar

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Menubar is a widget that represents a menubar.
type Menubar struct {
	MenuItems []MenuItem
}

// MenuItem is a single item in the menubar.
type MenuItem struct {
	Name      string
	clickable *widget.Clickable
	OnClick   func(gtx layout.Context)
}

// NewMenubar creates a new menubar.
func NewMenubar() *Menubar {
	return &Menubar{}
}

func (m *Menubar) AddMenuItem(name string, onClick func(gtx layout.Context)) {
	m.MenuItems = append(m.MenuItems, MenuItem{Name: name, clickable: new(widget.Clickable), OnClick: onClick})
}

func (menu Menubar) RenderMenu(gtx layout.Context, th *material.Theme) []layout.FlexChild {
	children := make([]layout.FlexChild, len(menu.MenuItems))
	for i, item := range menu.MenuItems {
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if item.clickable == nil {
				item.clickable = new(widget.Clickable)
			}

			if item.clickable.Clicked(gtx) {
				if item.OnClick != nil {
					item.OnClick(gtx)
				}
			}

			return material.Clickable(gtx, item.clickable, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Body1(th, item.Name).Layout(gtx)
				})
			})
		})
	}

	return children
}
