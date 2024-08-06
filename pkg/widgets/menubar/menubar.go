package menubar

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Menubar is a widget that represents a menubar.
type Menubar struct {
	MenuItems []*MenuItem
}

// MenuItem is a single item in the menubar.
type MenuItem struct {
	Clickable *widget.Clickable
	Flexed    bool
	OnClick   func(gtx layout.Context)
	Layout    func(gtx layout.Context, th *material.Theme) layout.Dimensions
}

func NewMenuItem(onClick func(gtx layout.Context), menuLayout func(gtx layout.Context, th *material.Theme) layout.Dimensions) *MenuItem {
	var clickable *widget.Clickable = nil
	if onClick != nil {
		clickable = new(widget.Clickable)
	}
	return &MenuItem{
		Clickable: clickable,
		OnClick:   onClick,
		Layout:    menuLayout,
	}
}

// NewMenubar creates a new menubar.
func NewMenubar() *Menubar {
	return &Menubar{}
}

func (m *Menubar) AddMenuItem(menuItem *MenuItem) {
	m.MenuItems = append(m.MenuItems, menuItem)
}

func (menu *Menubar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	children := make([]layout.FlexChild, len(menu.MenuItems))
	for i, item := range menu.MenuItems {
		layoutMenuItem := func(gtx layout.Context) layout.Dimensions {
			if item.Clickable != nil && item.OnClick != nil && item.Clickable.Clicked(gtx) {
				item.OnClick(gtx)
			}
			if item.Clickable != nil {
				return material.Clickable(gtx, item.Clickable, func(gtx layout.Context) layout.Dimensions {
					return item.Layout(gtx, th)
				})
			}
			return item.Layout(gtx, th)
		}
		if item.Flexed {
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
