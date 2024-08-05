package menubar

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/paint"
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
	Icon      *image.Image
	Scale     float32
	Updating  func(gtx layout.Context) string
	clickable *widget.Clickable
	OnClick   func(gtx layout.Context)
}

// NewMenubar creates a new menubar.
func NewMenubar() *Menubar {
	return &Menubar{}
}

func (m *Menubar) AddMenuItem(name string, onClick func(gtx layout.Context)) {
	m.MenuItems = append(m.MenuItems, MenuItem{
		Name:      name,
		clickable: new(widget.Clickable),
		OnClick:   onClick,
	})
}
func (m *Menubar) AddMenuIcon(icon *image.Image, scale float32, onClick func(gtx layout.Context)) {
	m.MenuItems = append(m.MenuItems, MenuItem{
		Name:      "",
		Icon:      icon,
		Scale:     scale,
		clickable: new(widget.Clickable),
		OnClick:   onClick,
	})
}
func (m *Menubar) AddMenuSeparator() {
	m.MenuItems = append(m.MenuItems, MenuItem{
		Name: "",
	})
}
func (m *Menubar) AddMenuLabel(text string, updating func(gtx layout.Context) string) {
	m.MenuItems = append(m.MenuItems, MenuItem{
		Name:     text,
		Updating: updating,
	})
}

func (menu *Menubar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	children := make([]layout.FlexChild, len(menu.MenuItems))
	for i, item := range menu.MenuItems {
		// update the item if it has an updating function
		if item.Updating != nil {
			item.Name = item.Updating(gtx)
		}

		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			itemLayout := layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if item.Icon != nil {
					icon := widget.Image{
						Src:   paint.NewImageOp(*item.Icon),
						Scale: item.Scale,
					}
					return icon.Layout(gtx)
				}
				return material.Body1(th, item.Name).Layout(gtx)
			})

			if item.clickable != nil {
				if item.clickable.Clicked(gtx) && item.OnClick != nil {
					item.OnClick(gtx)
				}
			} else {
				return itemLayout
			}

			return material.Clickable(gtx, item.clickable, func(gtx layout.Context) layout.Dimensions {
				return itemLayout
			})
		})
	}

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx, children...)
}
