package grid

import (
	"image"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Grid struct {
	Columns    int
	List       *widget.List
	ItemWidth  int
	ItemHeight int
}

type GridRowElement func(gtx layout.Context, index int) layout.Dimensions

func (g Grid) Layout(gtx layout.Context, theme *material.Theme, len int, r GridRowElement) layout.Dimensions {
	// create a 2D array of rows/cols because we need to render it in a list so its scrollable
	numberOfRows := len / g.Columns
	remainder := len % g.Columns
	if remainder > 0 {
		numberOfRows++
	}

	constraints := layout.Constraints{
		Min: image.Point{X: g.ItemWidth, Y: gtx.Constraints.Min.Y},
		Max: image.Point{X: g.ItemWidth, Y: g.ItemHeight},
		// Max: image.Point{X: g.ItemWidth, Y: gtx.Constraints.Max.Y},
	}
	originalConstraints := gtx.Constraints

	rows := make([][]layout.FlexChild, numberOfRows)
	for i := 0; i < numberOfRows; i++ {
		rows[i] = make([]layout.FlexChild, g.Columns)
		for j := 0; j < g.Columns; j++ {
			index := i*g.Columns + j
			rows[i][j] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if index >= len {
					return layout.Dimensions{}
				} else {
					gtx.Constraints = constraints
					return r(gtx, index)
				}
			})
		}
	}

	// reset constraints
	gtx.Constraints = originalConstraints

	return g.List.Layout(gtx, numberOfRows, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx, rows[index]...)
	})
}
