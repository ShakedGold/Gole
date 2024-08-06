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
	"github.com/ShakedGold/Gole/pkg/widgets/menubar/path"
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

	pathItem := path.NewPath(entries)

	// create the menu
	menu := menubar.NewMenubar()
	upMenuItem := &menubar.MenuItem{
		Clickable: new(widget.Clickable),
		OnClick: func(gtx layout.Context) {
			previousPath := filepath.Join(entries.Path, "..")
			previousEntries, err := entry.ReadPath(previousPath)
			if err != nil {
				log.Println(err)
				return
			}
			entries.Update(previousEntries)
		},
		Layout: func(gtx layout.Context, th *material.Theme) layout.Dimensions {
			return widget.Image{
				Src: paint.NewImageOp(*up),
			}.Layout(gtx)
		},
	}

	menu.AddMenuItem(upMenuItem)
	menu.AddMenuItem(pathItem.MenuItem)

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
				entrys, err := entry.ReadPath(pathItem.PathEditor.Text())
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
					}

					return layoutEntries
				}),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
