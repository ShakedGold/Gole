package clickable

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ClickableMenuItem struct {
	Clickable      *widget.Clickable
	Flexed         bool
	OnClick        func(gtx layout.Context)
	LayoutCallback func(gtx layout.Context, th *material.Theme) layout.Dimensions
}

func NewClickableMenuItem(onClick func(gtx layout.Context), layoutCallback func(gtx layout.Context, th *material.Theme) layout.Dimensions) ClickableMenuItem {
	return ClickableMenuItem{
		Clickable:      new(widget.Clickable),
		OnClick:        onClick,
		LayoutCallback: layoutCallback,
	}
}

func (c ClickableMenuItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for c.Clickable != nil && c.Clickable.Clicked(gtx) && c.OnClick != nil {
		c.OnClick(gtx)
	}
	return material.Clickable(gtx, c.Clickable, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return c.LayoutCallback(gtx, th)
		})
	})
}

func (c ClickableMenuItem) IsFlexed() bool {
	return c.Flexed
}
