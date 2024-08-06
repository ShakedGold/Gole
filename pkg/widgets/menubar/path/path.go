package path

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ShakedGold/Gole/pkg/widgets/entry"
	"github.com/ShakedGold/Gole/pkg/widgets/menubar"
)

type Path struct {
	*menubar.MenuItem
	PathEditor *widget.Editor
	Entries    *entry.Entries
}

func NewPath(entries *entry.Entries) *Path {
	path := &Path{
		Entries: entries,
	}
	path.PathEditor = new(widget.Editor)
	path.PathEditor.SingleLine = true
	path.PathEditor.Submit = true
	path.PathEditor.SetText(entries.Path)
	path.MenuItem = menubar.NewMenuItem(nil, func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		return material.Editor(th, path.PathEditor, "Path").Layout(gtx)
	})
	path.MenuItem.Flexed = true
	return path
}

func (p *Path) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return p.MenuItem.Layout(gtx, th)
}
