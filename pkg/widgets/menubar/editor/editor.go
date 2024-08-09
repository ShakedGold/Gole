package editor

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type EditorInputItem struct {
	Editor         *widget.Editor
	Flexed         bool
	UpdateCallback func(*widget.Editor)
}

func (e EditorInputItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if e.UpdateCallback != nil {
		e.UpdateCallback(e.Editor)
	}
	return material.Editor(th, e.Editor, "Path").Layout(gtx)
}

func (e EditorInputItem) IsFlexed() bool {
	return e.Flexed
}
