package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ShakedGold/Gole/pkg/assets"
	"github.com/ShakedGold/Gole/pkg/explorer"
	"github.com/ShakedGold/Gole/pkg/menubar"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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
	entries = entries.Prepare()

	// watch the home directory
	watcher, err := explorer.Watcher(entries.Path)
	if err != nil {
		return err
	}

	// watch for events, if there is a change in the filesystem
	go func() {
		explorer.Watch(watcher, &entries)
	}()

	// Create a List widget
	var list widget.List
	list.Axis = layout.Vertical

	// create the menu
	menu := menubar.NewMenubar()
	menu.AddMenuItem("File", func(gtx layout.Context) {
		log.Println("File clicked")
	})
	menu.AddMenuItem("Edit", func(gtx layout.Context) {
		log.Println("Edit clicked")
	})

	// Create a clickable for each entry
	// clickables := make([]widget.Clickable, len(entries.Entries))

	var ops op.Ops

	folderIcon, err := assets.GetImage("folder.png")
	if err != nil {
		return err
	}
	fileIcon, err := assets.GetImage("file.png")
	if err != nil {
		return err
	}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			var newEntries explorer.Entries

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					items := menu.RenderMenu(gtx, theme)

					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx, items...)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return list.Layout(gtx, len(entries.Entries), func(gtx layout.Context, index int) layout.Dimensions {
						clickable := entries.Entries[index].Clickable

						if clickable == nil {
							clickable = new(widget.Clickable)
							entries.Entries[index].Clickable = clickable
						}

						if clickable.Clicked(gtx) {
							entry := entries.Entries[index]
							newEntries, err = entry.EntryAction(watcher)
							if err != nil {
								log.Println(err)
							}
						}

						return material.Clickable(gtx, clickable, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Horizontal,
								}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{
											Right: unit.Dp(8),
											Left:  unit.Dp(8),
										}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											var icon widget.Image

											if entries.Entries[index].IsFolder {
												icon = widget.Image{
													Src:   paint.NewImageOp(folderIcon),
													Scale: 0.1,
												}
											} else {
												icon = widget.Image{
													Src:   paint.NewImageOp(fileIcon),
													Scale: 0.1,
												}
											}

											return icon.Layout(gtx)
										})
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										// get last part of path
										var name string
										if entries.Entries[index].Alias != "" {
											name = entries.Entries[index].Alias
										} else {
											name = filepath.Base(entries.Entries[index].Path)
										}
										return material.Body1(theme, name).Layout(gtx)
									}),
								)
							})
						})
					})
				}),
			)

			menu.RenderMenu(gtx, theme)

			// switch if new entries are available
			if len(newEntries.Entries) > 0 {
				entries = newEntries
			}

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
