package main

import (
	"log"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ShakedGold/Gole/pkg/assets"
	"github.com/ShakedGold/Gole/pkg/explorer"
	"github.com/ShakedGold/Gole/pkg/widgets/entry"
	"github.com/ShakedGold/Gole/pkg/widgets/menubar"
	"github.com/ShakedGold/Gole/pkg/widgets/menubar/clickable"
	"github.com/ShakedGold/Gole/pkg/widgets/menubar/editor"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("Gole"))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()

	entries, err := explorer.Home()
	if err != nil {
		return err
	}
	entries, err = entries.Prepare()
	if err != nil {
		return err
	}

	// watch the home directory
	watcher, err := explorer.Watcher(entries.Path)
	if err != nil {
		return err
	}

	// watch for events, if there is a change in the filesystem
	go func() {
		explorer.Watch(watcher, entries)
	}()

	// get up asset
	up, err := assets.GetImage("up.png")
	if err != nil {
		return err
	}

	list, err := assets.GetImage("list.png")
	if err != nil {
		return err
	}

	grid, err := assets.GetImage("grid.png")
	if err != nil {
		return err
	}

	pathEditor := new(widget.Editor)
	pathEditor.SetText(entries.Path)

	// create the menu
	menu := menubar.NewMenubar()
	upMenuItem := clickable.ClickableMenuItem{
		Clickable: new(widget.Clickable),
		OnClick: func(gtx layout.Context) {
			previousPath := filepath.Join(entries.Path, "..")
			previousEntries, err := entry.ReadPath(previousPath)
			if err != nil {
				log.Println(err)
				return
			}
			entries.Update(previousEntries)
			pathEditor.SetText(previousEntries.Path)
		},
		LayoutCallback: func(gtx layout.Context, th *material.Theme) layout.Dimensions {
			return widget.Image{
				Src:   paint.NewImageOp(*up),
				Scale: 0.5,
			}.Layout(gtx)
		},
	}

	pathMenuItem := editor.EditorInputItem{
		Editor: pathEditor,
		Flexed: true,
	}

	viewMenuItem := clickable.ClickableMenuItem{
		Clickable: new(widget.Clickable),
		OnClick: func(gtx layout.Context) {
			entries.ViewMode = 1 - entries.ViewMode
		},
		LayoutCallback: func(gtx layout.Context, th *material.Theme) layout.Dimensions {
			var image widget.Image
			var label material.LabelStyle

			if entries.ViewMode == entry.ViewModeGrid {
				image = widget.Image{
					Src:   paint.NewImageOp(*grid),
					Scale: 0.6,
				}
				label = material.H6(th, "Grid")
			} else {
				image = widget.Image{
					Src:   paint.NewImageOp(*list),
					Scale: 0.5,
				}
				label = material.H6(th, "List")
			}
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(image.Layout),
				layout.Rigid(label.Layout),
			)
		},
	}

	menu.AddMenuItem(upMenuItem)
	menu.AddMenuItem(pathMenuItem)
	menu.AddMenuItem(viewMenuItem)

	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Process events that arrived between the last frame and this one.
			for {
				// wait for any keyboard input
				_, ok := e.Source.Event(
					key.Filter{
						Name: key.NameReturn,
					},
				)
				if !ok {
					break
				}

				// if the path editor is focused, update the path
				entrys, err := entry.ReadPath(pathEditor.Text())
				if err != nil {
					log.Println(err)
					break
				}
				entries.Update(entrys)
			}

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return menu.Layout(gtx, theme)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					layoutEntries, updatedEntries, err := entries.Layout(gtx, theme, watcher)
					if err != nil {
						log.Println(err)
						return layout.Dimensions{}
					}

					if updatedEntries != nil {
						entries.Update(updatedEntries)
						pathEditor.SetText(updatedEntries.Path)
					}

					return layoutEntries
				}),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
