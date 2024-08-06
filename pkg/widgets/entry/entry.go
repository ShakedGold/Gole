package entry

import (
	"image"
	"os"
	"path/filepath"
	"sort"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ShakedGold/Gole/pkg/assets"
	"github.com/ShakedGold/Gole/pkg/widgets"
	"github.com/ShakedGold/Gole/pkg/widgets/grid"
	"github.com/fsnotify/fsnotify"
	"github.com/skratchdot/open-golang/open"
)

type Entry struct {
	Path      string
	Alias     string
	Width     int
	Height    int
	IsFolder  bool
	Clickable *widget.Clickable
	Icon      *image.Image
}

type Entries struct {
	Entries []Entry
	Path    string
	List    *widget.List
	Grid    *grid.Grid
}

func CreateFile(path string, alias string) (Entry, error) {
	icon, err := assets.GetImage("file.png")
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Path:      path,
		Alias:     alias,
		IsFolder:  false,
		Clickable: new(widget.Clickable),
		Icon:      icon,
		Width:     200,
		Height:    200,
	}, nil
}
func CreateFolder(path string, alias string) (Entry, error) {
	icon, err := assets.GetImage("folder.png")
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Path:      path,
		Alias:     alias,
		IsFolder:  true,
		Clickable: new(widget.Clickable),
		Icon:      icon,
		Width:     200,
		Height:    200,
	}, nil
}

func ReadPath(path string) (*Entries, error) {
	osEntries, err := os.ReadDir(path)

	if err != nil {
		return &Entries{}, err
	}

	list := new(widget.List)
	list.Axis = layout.Vertical

	grid := widgets.Grid(9)

	entries := Entries{
		Entries: []Entry{},
		Path:    path,
		List:    list,
		Grid:    &grid,
	}

	for _, osEntry := range osEntries {
		var entry Entry
		if osEntry.IsDir() {
			entry, err = CreateFolder(filepath.Join(path, osEntry.Name()), osEntry.Name())
			if err != nil {
				return &Entries{}, err
			}
		} else {
			entry, err = CreateFile(filepath.Join(path, osEntry.Name()), osEntry.Name())
			if err != nil {
				return &Entries{}, err
			}
		}

		entries.Entries = append(entries.Entries, entry)
	}

	return &entries, nil
}

func (e *Entry) Action(watcher *fsnotify.Watcher) (*Entries, error) {
	if e.IsFolder {
		entries, err := ReadPath(e.Path)
		if err != nil {
			return &Entries{}, err
		}

		entries, err = entries.Prepare()
		if err != nil {
			return &Entries{}, err
		}

		// remove from watch
		watcher.Remove(e.Path)

		// watch new folder
		watcher.Add(entries.Path)

		return entries, nil
	} else {
		open.Run(e.Path)
	}

	return &Entries{}, nil
}

func (entries *Entries) Prepare() (*Entries, error) {
	// sort entries
	sort.Slice(entries.Entries, ByIsDir(entries.Entries))

	return entries, nil
}

func ByIsDir(entries []Entry) func(a, b int) bool {
	return func(a, b int) bool {
		// sort by IsDir, then by Name
		if entries[a].IsFolder == entries[b].IsFolder {
			return entries[a].Path < entries[b].Path
		}
		return entries[a].IsFolder
	}
}

func (e Entry) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	constraints := layout.Constraints{
		Min: image.Point{X: e.Width, Y: e.Height},
		Max: image.Point{X: e.Width, Y: e.Height},
	}
	originalConstraints := gtx.Constraints
	gtx.Constraints = constraints

	clickable := material.Clickable(gtx, e.Clickable, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{
				Alignment: layout.Center,
			}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if e.Icon == nil {
								return layout.Dimensions{}
							}
							icon := widget.Image{
								Src: paint.NewImageOp(*e.Icon),
							}

							return icon.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top: unit.Dp(10),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Dimensions{}
							})
						}), // spacer
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// get last part of path
							var name string
							if e.Alias != "" {
								name = e.Alias
							} else {
								name = filepath.Base(e.Path)
							}

							return material.Body1(theme, name).Layout(gtx)
						}))
				}),
			)
		})
	})

	gtx.Constraints = originalConstraints

	return clickable
}

func (entries *Entries) Layout(gtx layout.Context, theme *material.Theme, watcher *fsnotify.Watcher) (layout.Dimensions, *Entries, error) {
	var layoutErr error
	var updatedEntries *Entries

	// get calculated width
	width := gtx.Constraints.Max.X
	// height := gtx.Constraints.Max.Y

	if len(entries.Entries) == 0 {
		return layout.Dimensions{}, nil, nil
	}

	// calculate number of columns
	columns := width / entries.Entries[0].Width
	if columns == 0 {
		columns = 1
	}

	// update grid columns
	entries.Grid.Columns = columns

	layout := entries.Grid.Layout(gtx, theme, len(entries.Entries), func(gtx layout.Context, index int) layout.Dimensions {
		clickable := entries.Entries[index].Clickable

		if clickable == nil {
			entries.Entries[index].Clickable = new(widget.Clickable)
		}

		if clickable.Clicked(gtx) {
			entry := entries.Entries[index]
			newEntries, err := entry.Action(watcher)
			if err != nil {
				layoutErr = err
			}

			// switch if new entries are available
			if len(newEntries.Entries) > 0 {
				updatedEntries = newEntries
			}
		}

		return entries.Entries[index].Layout(gtx, theme)
	})

	if layoutErr != nil {
		return layout, nil, layoutErr
	}

	return layout, updatedEntries, nil
}
