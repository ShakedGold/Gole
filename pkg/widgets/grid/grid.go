package grid

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Grid struct {
	Columns int
	List    *widget.List
}

type GridRowElement func(gtx layout.Context, index int) layout.Dimensions

func (g Grid) Layout(gtx layout.Context, theme *material.Theme, len int, r GridRowElement) layout.Dimensions {
	// create a 2D array of rows/cols because we need to render it in a list so its scrollable
	numberOfRows := len / g.Columns
	remainder := len % g.Columns
	if remainder > 0 {
		numberOfRows++
	}

	rows := make([][]layout.FlexChild, numberOfRows)
	for i := 0; i < numberOfRows; i++ {
		rows[i] = make([]layout.FlexChild, g.Columns)
		for j := 0; j < g.Columns; j++ {
			index := i*g.Columns + j
			rows[i][j] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if index >= len {
					return layout.Dimensions{}
				} else {
					return r(gtx, index)
				}
			})
		}
	}

	return g.List.Layout(gtx, numberOfRows, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx, rows[index]...)
	})
}
